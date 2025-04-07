package hasher

import "golang.org/x/crypto/bcrypt"

//go:generate go tool mockgen -destination mock_$GOFILE -package=$GOPACKAGE . PasswordHasher
type PasswordHasher interface {
	Hash(password string) (string, error)
	Match(password, hashedPassword string) bool
}

type bcryptHasher struct{}

func NewBcrypt() PasswordHasher {
	return &bcryptHasher{}
}

func (h *bcryptHasher) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (h *bcryptHasher) Match(password, hashedPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}

	return true
}
