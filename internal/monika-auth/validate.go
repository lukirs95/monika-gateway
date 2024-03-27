package monikaauth

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	sdk "github.com/lukirs95/monika-gosdk/pkg/types"
)

type Session struct {
	UserId int64
	Role   sdk.UserRole
}

func GetSessionFrom(r *http.Request) (*Session, error) {
	ctx := r.Context()
	session := ctx.Value(MonikaSession)
	if session == nil {
		return nil, fmt.Errorf("no valid session")
	}
	return session.(*Session), nil
}

func (auth *MonikaAuth) ValidateJWT(tokenString string) (*Session, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt Token parsing error")
		}
		return auth.params.sessionSalt, nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*jwt.StandardClaims); ok {
		if claims.ExpiresAt < time.Now().Unix() && claims.ExpiresAt != 0 {
			return nil, fmt.Errorf("token expired")
		}

		userId, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("malicious user id")
		}
		session := &Session{
			UserId: userId,
			Role:   sdk.UserRole(claims.Audience),
		}

		return session, nil
	} else {
		return nil, fmt.Errorf("can not parse JWT Token")
	}
}
