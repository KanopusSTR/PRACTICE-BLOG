package pkg

import "crypto/sha256"

func GetStringHash(password string) []byte {
	hash := sha256.New()
	hash.Write([]byte(password))
	sum := hash.Sum(nil)
	return sum
}
