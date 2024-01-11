SHELL=/bin/bash -o pipefail

PRODUCT_OWNER_NAME := praromvik
PRODUCT_NAME       := praromvik
GO_PKG   := github.com/praromvik
BIN      := praromvik


OS   := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))

BUILD_DIRS  := bin/$(OS)_$(ARCH)     \
               .go/bin/$(OS)_$(ARCH) \
               .go/cache             \
               hack/config


$(BUILD_DIRS):
	@mkdir -p $@

GO_VERSION       ?= 1.21
BUILD_IMAGE      ?= ghcr.io/appscode/golang-dev:$(GO_VERSION)


SRC_PKGS := pkg cmd handlers
SRC_DIRS := $(SRC_PKGS)



fmt: $(BUILD_DIRS)
	@docker run                                                 \
	    -i                                                      \
	    --rm                                                    \
	    -u $$(id -u):$$(id -g)                                  \
	    -v $$(pwd):/src                                         \
	    -w /src                                                 \
	    -v $$(pwd)/.go/cache:/.cache                            \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin                \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH)  \
	    $(BUILD_IMAGE)                                          \
	    /bin/bash -c "                                          \
	        REPO_PKG=$(GO_PKG)                                  \
	        ./hack/fmt.sh $(SRC_DIRS)                           \
	    "



DOCKER_REPO_ROOT := /go/src/$(GO_PKG)/$(REPO)

add-license:
	@echo "Adding license header"
	@docker run --rm 	                                 \
		-u $$(id -u):$$(id -g)                           \
		-v /tmp:/.cache                                  \
		-v $$(pwd):$(DOCKER_REPO_ROOT)                   \
		-w $(DOCKER_REPO_ROOT)                           \
		$(BUILD_IMAGE)                                   \
		ltag -t "./hack/license" --excludes "vendor contrib third_party libbuild" -v