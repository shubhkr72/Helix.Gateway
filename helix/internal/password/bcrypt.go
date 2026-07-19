package password

import "golang.org/x/crypto/bcrypt"

// HashPassword converts a plain text password into a bcrypt hash.
func HashPassword(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// VerifyPassword compares a plain password with its bcrypt hash.
func VerifyPassword(password, hash string) error {

	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}