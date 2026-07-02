package seguranca

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {}

func NovoBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	bytes, erro := bcrypt.GenerateFromPassword([]byte(password), 12)
	if erro != nil {
		return "", erro
	}
	return string(bytes), nil
}

func (h *BcryptHasher) Compare(hash, password string) bool {
	erro := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return erro == nil
}
