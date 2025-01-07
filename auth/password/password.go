package password

import "golang.org/x/crypto/bcrypt"

var defaultStrength = bcrypt.DefaultCost

func GeneratePasswordHash(password string) []byte {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), defaultStrength)
	if err != nil {
		panic(err)
	}
	return passwordHash
}

func CompareHashPassword(passwordHash, password []byte) bool {
	return bcrypt.CompareHashAndPassword(passwordHash, password) == nil
}
