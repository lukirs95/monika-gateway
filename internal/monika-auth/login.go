package monikaauth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	sdk "github.com/lukirs95/monika-gosdk/pkg/types"
)

func (auth *MonikaAuth) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	var askingUser sdk.Auth
	if err := json.NewDecoder(r.Body).Decode(&askingUser); err != nil {
		log.Print(err)
		http.Error(w, "could not parse request", http.StatusBadRequest)
		return
	}

	jwttoken, err := auth.Login(askingUser)
	if err != nil {
		log.Print(err)
		http.Error(w, "wrong username or password", http.StatusForbidden)
	}

	cookieAge := 0
	if askingUser.Expires {
		cookieAge = int(auth.params.sessionExpiration.Seconds())
	}

	jwt := &http.Cookie{
		Name:     "Authorization",
		Value:    jwttoken,
		MaxAge:   cookieAge,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, jwt)
}

func (auth *MonikaAuth) Login(askingUser sdk.Auth) (string, error) {
	askingUser.Username.Sanitize()
	// check if user is in database
	user, err := auth.userDB.GetUserByUsername(askingUser.Username)
	if err != nil {
		return "", err
	}

	// hash password and compare with hashed password in database
	isValid, err := auth.comparePassword(askingUser.Password, user.Password)
	if err != nil {
		return "", err
	}

	if !isValid {
		return "", fmt.Errorf("wrong password")
	}

	var expiresAt int64 = 0
	if askingUser.Expires {
		expiresAt = time.Now().Add(auth.params.sessionExpiration).Unix()
	}

	claims := jwt.StandardClaims{
		Subject:   fmt.Sprintf("%d", user.UserId),
		Audience:  string(user.Role),
		ExpiresAt: expiresAt,
	}

	// generate JWToken
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(auth.params.sessionSalt)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
