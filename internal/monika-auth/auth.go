package monikaauth

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/alexedwards/argon2id"
	monikadb "github.com/lukirs95/monika-gateway/internal/monika-db/dbrepo"
)

type MonikaAuth struct {
	protectedRoutes http.Handler
	userDB          monikadb.UserDatabaseRepo
	params          MonikaAuthParams
}

func NewMonikaAuth(db monikadb.UserDatabaseRepo, params MonikaAuthParams) *MonikaAuth {
	return &MonikaAuth{
		userDB: db,
		params: params,
	}
}

type MonikaAuthParams struct {
	sessionExpiration time.Duration
	sessionSalt       []byte
	argon2Params      argon2id.Params
}

func NewMonikaAuthParams(expire time.Duration, sessionSalt []byte) MonikaAuthParams {
	return MonikaAuthParams{
		sessionExpiration: expire,
		sessionSalt:       sessionSalt,
		argon2Params: argon2id.Params{
			Memory:      argon2id.DefaultParams.Memory,
			Iterations:  argon2id.DefaultParams.Iterations,
			Parallelism: uint8(runtime.NumCPU()),
			SaltLength:  argon2id.DefaultParams.SaltLength,
			KeyLength:   argon2id.DefaultParams.KeyLength,
		},
	}
}

type MonikaSessionKey string

const MonikaSession MonikaSessionKey = "MonikaSession"

func (auth *MonikaAuth) SecureHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/authenticate" {
			auth.protectedRoutes.ServeHTTP(w, r)
		}

		token, err := r.Cookie("Authorization")
		if err != nil {
			http.Error(w, "no bearer token provided", http.StatusUnauthorized)
			return
		}

		session, err := auth.ValidateJWT(token.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctxWithSession := context.WithValue(r.Context(), MonikaSession, session)
		rWithCtxSession := r.WithContext(ctxWithSession)
		next.ServeHTTP(w, rWithCtxSession)
	}
}

func (auth *MonikaAuth) Secure(handle http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/authenticate" {
			auth.protectedRoutes.ServeHTTP(w, r)
		}

		token, err := r.Cookie("Authorization")
		if err != nil {
			http.Error(w, "no token provided", http.StatusUnauthorized)
			return
		}

		session, err := auth.ValidateJWT(token.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctxWithSession := context.WithValue(r.Context(), MonikaSession, session)
		rWithCtxSession := r.WithContext(ctxWithSession)
		handle(w, rWithCtxSession)
	}
}
