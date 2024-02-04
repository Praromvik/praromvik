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

.PHONY: ci
ci: verify check-license # lint build

.PHONY: verify
verify: verify-gen verify-modules

.PHONY: verify-gen
verify-gen: gen fmt
	@if !(git diff --exit-code HEAD); then \
		echo "files are out of date, run make gen fmt"; exit 1; \
	fi

gen:
	@true

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

ADDTL_LINTERS   := gofmt,goimports,unparam

.PHONY: lint
lint: $(BUILD_DIRS)
	@echo "running linter"
	@docker run                                                 \
	    -i                                                      \
	    --rm                                                    \
	    -u $$(id -u):$$(id -g)                                  \
	    -v $$(pwd):/src                                         \
	    -w /src                                                 \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin                \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH)  \
	    -v $$(pwd)/.go/cache:/.cache                            \
	    --env GO111MODULE=on                                    \
	    --env GOFLAGS="-mod=vendor"                             \
	    $(BUILD_IMAGE)                                          \
	    golangci-lint run --enable $(ADDTL_LINTERS) --timeout=10m --skip-files="generated.*\.go$\" --skip-dirs-use-default --skip-dirs=client,vendor

.PHONY: verify-modules
verify-modules:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor
	@if !(git diff --exit-code HEAD); then \
		echo "go module files are out of date"; exit 1; \
	fi

DOCKER_REPO_ROOT := /go/src/$(GO_PKG)/$(REPO)

.PHONY: add-license
add-license:
	@echo "Adding license header"
	@docker run --rm 	                                 \
		-u $$(id -u):$$(id -g)                           \
		-v /tmp:/.cache                                  \
		-v $$(pwd):$(DOCKER_REPO_ROOT)                   \
		-w $(DOCKER_REPO_ROOT)                           \
		$(BUILD_IMAGE)                                   \
		ltag -t "./hack/license" --excludes "vendor contrib third_party libbuild" -v

.PHONY: check-license
check-license:
	@echo "Checking files for license header"
	@docker run --rm 	                                 \
		-u $$(id -u):$$(id -g)                           \
		-v /tmp:/.cache                                  \
		-v $$(pwd):$(DOCKER_REPO_ROOT)                   \
		-w $(DOCKER_REPO_ROOT)                           \
		$(BUILD_IMAGE)                                   \
		ltag -t "./hack/license" --excludes "vendor contrib third_party libbuild" --check -v

.PHONY: clean
clean:
	rm -rf .go bin