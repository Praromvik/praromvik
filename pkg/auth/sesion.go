/*
MIT License

Copyright (c) 2024 Praromvik

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package auth

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	math "math/rand"
	"net/http"
	"os"

	"github.com/praromvik/praromvik/models"
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
	sessionTokenName = randStr(math.Intn(11) + 10)
	redisStore.KeyPrefix(os.Getenv("SESSION_KEY"))
}

func randStr(length int) string {
	randomBytes := make([]byte, length)
	return hex.EncodeToString(randomBytes)[:length]
}

func StoreAuthenticated(w http.ResponseWriter, r *http.Request, u *user.User, v bool) error {
	session, err := redisStore.Get(r, sessionTokenName)
	if err != nil {
		return err
	}
	session.Values[authenticated] = true
	if u != nil {
		session.Values[models.Role] = u.Role
	}
	if !v {
		session.Options.MaxAge = -1
	}
	return session.Save(r, w)
}

func IsAuthenticated(r *http.Request) (bool, error) {
	session, err := redisStore.Get(r, sessionTokenName)
	if err != nil {
		return false, err
	}
	authValue, _ := session.Values[authenticated].(bool)
	return authValue, nil
}

func GetSessionRole(r *http.Request) (models.RoleType, error) {
	session, err := redisStore.Get(r, sessionTokenName)
	if err != nil {
		return "", err
	}
	role, ok := session.Values[models.Role]
	if !ok || role == nil {
		return models.RoleType(""), nil
	}

	switch role.(string) {
	case "Admin":
		return models.Admin, nil
	case "Moderator":
		return models.Moderator, nil
	case "Student":
		return models.Student, nil
	case "Trainer":
		return models.Trainer, nil
	default:
		return "", fmt.Errorf("unknown role type: %s", role.(string))
	}
}

func SessionValid(r *http.Request) (bool, error) {
	session, err := redisStore.Get(r, sessionTokenName)
	if err != nil {
		return false, err
	}
	return !session.IsNew, nil
}
