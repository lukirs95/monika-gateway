package monikaauth

import (
	"runtime"

	"github.com/alexedwards/argon2id"
	sdk "github.com/lukirs95/monika-gosdk/pkg/types"
)

func (auth *MonikaAuth) hashPassword(user *sdk.User) error {
	hash, err := argon2id.CreateHash(string(user.Password), &auth.params.argon2Params)
	user.Password = sdk.Password(hash)
	return err
}

func (auth *MonikaAuth) comparePassword(password sdk.Password, hashedPassword sdk.Password) (bool, error) {
	return argon2id.ComparePasswordAndHash(string(password), string(hashedPassword))
}

func HashDefaultPassword(password string) (string, error) {
	return argon2id.CreateHash(password, &argon2id.Params{
		Memory:      argon2id.DefaultParams.Memory,
		Iterations:  argon2id.DefaultParams.Iterations,
		Parallelism: uint8(runtime.NumCPU()),
		SaltLength:  argon2id.DefaultParams.SaltLength,
		KeyLength:   argon2id.DefaultParams.KeyLength,
	})
}
