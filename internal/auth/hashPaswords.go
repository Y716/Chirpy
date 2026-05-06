package auth

import "github.com/alexedwards/argon2id"

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func ComparePasswordHash(password, hash string) (bool, error) {
	flag, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}

	return flag, nil
}
