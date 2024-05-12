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
	"github.com/praromvik/praromvik/models/utils"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/praromvik/praromvik/models/db/client"
	"github.com/praromvik/praromvik/models/user"

	rstore "github.com/rbcervilla/redisstore/v8"
)

var (
	redisStore       *rstore.RedisStore
	sessionTokenName = "PRAROMVIK"
)

func init() {
	var err error
	redisStore, err = rstore.NewRedisStore(context.Background(), client.Redis)
	if err != nil {
		log.Fatal("failed to create redis store: ", err)
	}
	redisStore.KeyPrefix(os.Getenv(utils.SessionKey))
}

func StoreAuthenticated(w http.ResponseWriter, r *http.Request, u *user.User, valid bool) error {
	session, err := redisStore.Get(r, sessionTokenName)
	if err != nil {
		return err
	}
	session.Values[utils.Authenticated] = true
	if u != nil {
		session.Values[utils.UUID] = u.UUID
		session.Values[utils.Role] = u.Role
		session.Values[utils.UserName] = u.UserName
	}
	session.Values[utils.UserIP] = getIpAddress(r)
	session.Values[utils.UserAgent] = r.UserAgent()
	if !valid {
		session.Options.MaxAge = -1
	}
	return session.Save(r, w)
}

func IsAuthenticated(r *http.Request) (bool, error) {
	session, err := redisStore.Get(r, sessionTokenName)
	if err != nil {
		return false, err
	}
	authValue := session.Values[utils.Authenticated].(bool)
	return authValue, nil
}

func GetSessionRole(r *http.Request) (utils.RoleType, error) {
	session, err := redisStore.Get(r, sessionTokenName)
	if err != nil {
		return "", err
	}
	role, ok := session.Values[utils.Role]
	if !ok || role == nil {
		return utils.None, nil
	}
	return utils.RoleType(role.(string)), nil
}

func SessionValid(r *http.Request) (bool, error) {
	session, err := redisStore.Get(r, sessionTokenName)
	if err != nil {
		return false, err
	}
	return !session.IsNew &&
		session.Values[utils.UserIP] == getIpAddress(r) &&
		session.Values[utils.UserAgent] == r.UserAgent(), nil
}

func GetUserInfoFromSession(r *http.Request) (*utils.Info, error) {
	session, err := redisStore.Get(r, sessionTokenName)
	if err != nil {
		return nil, err
	}
	return &utils.Info{Name: session.Values[utils.UserName].(string), UUID: session.Values[utils.UUID].(string)}, nil
}

func getIpAddress(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}
