package password

import "golang.org/x/crypto/bcrypt"

var defaultStrength = bcrypt.DefaultCost

// CreatePassword returns a hashed version of the given password.
func CreatePassword(pw string) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), defaultStrength)
	if err != nil {
		panic(err)
	}
	return hashedPassword
}

// ComparePassword compares a hashed password with its possible plaintext equivalent.
func ComparePassword(hashedPassword, password []byte) bool {
	return bcrypt.CompareHashAndPassword(hashedPassword, password) == nil
}
