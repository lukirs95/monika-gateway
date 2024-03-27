package monikaauth

import (
	"fmt"
	"net/http"
	"strings"
)

func (auth *MonikaAuth) HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		http.Error(w, "no token in header", http.StatusBadRequest)
		return
	}

	token := strings.Split(cookie.Value, "Bearer ")
	if len(token) != 2 {
		http.Error(w, "wrong token format", http.StatusBadRequest)
		return
	}

	jwt := &http.Cookie{
		Name:     "Authorization",
		Value:    fmt.Sprintf("Bearer %s", token[1]),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, jwt)
	w.WriteHeader(http.StatusOK)
}
