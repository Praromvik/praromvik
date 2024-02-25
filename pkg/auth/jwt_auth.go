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
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/praromvik/praromvik/models/user"
	perror "github.com/praromvik/praromvik/pkg/error"

	"github.com/golang-jwt/jwt"
)

var (
	cookieName = "JWT"
	jwtSecret  = "JWT_SECRET"
)

type jwtClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWTAndSetCookie(w http.ResponseWriter, user *user.User) error {
	str, err := generateJWT(user)
	if err != nil {
		return err
	}
	// set token on cookies storage
	http.SetCookie(w, &http.Cookie{Name: cookieName, Value: str})
	return nil
}

func generateJWT(user *user.User) (tokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &jwtClaim{
		Email:    user.Email,
		Username: user.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := []byte(os.Getenv(jwtSecret))
	tokenString, err = token.SignedString(jwtKey)
	return
}

func VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie(cookieName)
		if err != nil {
			perror.HandleError(writer, http.StatusUnauthorized, "", err)
			return
		}
		if err := validateToken(cookie.Value); err != nil {
			perror.HandleError(writer, http.StatusUnauthorized, "failed to validate jwt token", err)
		} else {
			next.ServeHTTP(writer, request)
		}
	})
}

func validateToken(signedToken string) error {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			jwtKey := []byte(os.Getenv(jwtSecret))
			return jwtKey, nil
		},
	)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(*jwtClaim)
	if !ok {
		return fmt.Errorf("couldn't parse claims")
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return fmt.Errorf("token expired")
	}
	return nil
}

func UnsetJWTInCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		fmt.Println(err)
		return
	}
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)
}
