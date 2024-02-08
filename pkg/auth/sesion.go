package auth

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/praromvik/praromvik/models/client"
	"github.com/praromvik/praromvik/models/user"

	rstore "github.com/rbcervilla/redisstore/v8"
)

var (
	redisStore       *rstore.RedisStore
	sessionTokenName string
	authenticated    = "authenticated"
)

func init() {
	var err error
	redisStore, err = rstore.NewRedisStore(context.Background(), client.Redis)
	if err != nil {
		log.Fatal("failed to create redis store: ", err)
	}
	sessionTokenName = randStr(rand.Intn(11) + 10)
	redisStore.KeyPrefix(os.Getenv("SESSION_KEY"))
}

func randStr(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	// Read b number of numbers
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func StoreAuthenticated(w http.ResponseWriter, r *http.Request, u *user.User, v bool) error {
	session, err := redisStore.Get(r, sessionTokenName)
	if err != nil {
		return err
	}
	session.Values[authenticated] = v
	if u != nil {
		session.Values["role"] = u.Role
	}
	return session.Save(r, w)
}

func IsAuthenticated(r *http.Request) (bool, error) {
	session, err := redisStore.Get(r, sessionTokenName)
	if err != nil {
		return false, err
	}
	a, _ := session.Values[authenticated]
	if a == nil {
		return false, nil
	}
	return a.(bool), nil
}

func IsAdminOrModeratorAuthenticated(r *http.Request) (bool, error) {
	session, err := redisStore.Get(r, sessionTokenName)
	if err != nil {
		return false, err
	}
	auth, _ := session.Values[authenticated]
	role, _ := session.Values["role"]
	if auth == nil || role == nil {
		return false, nil
	}
	return role.(string) == "admin", nil
}

func SessionValid(r *http.Request) (bool, error) {
	session, err := redisStore.Get(r, sessionTokenName)
	if err != nil {
		return false, err
	}
	return !session.IsNew, nil
}
