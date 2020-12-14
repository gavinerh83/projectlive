//Package secure contains functions that provide security to overall application.
//It contains functions that generate and validate JWT for each session login.
package secure

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

//MyClaims is the claims struct that is used in the generation of JWT
type MyClaims struct {
	jwt.StandardClaims
	SessionID string
}

var key []byte
var encryptionkey string

func getEnv(k string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return (os.Getenv(k))
}

//GenerateJWT creates a JWT for each user session login
func GenerateJWT(c *MyClaims) (string, error) {
	key = []byte(getEnv("JWT_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("Error in createToken when signing token")
	}
	return signedToken, nil
}

//ParseToken takes in a signedToken and returns the custom claims object and any validation errors
func ParseToken(signedToken string) (*MyClaims, error) {
	claims := &MyClaims{}
	t, err := jwt.ParseWithClaims(signedToken, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS512.Alg() { //the token here is not yet verified
			return nil, fmt.Errorf("Invalid signing algorithm")
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("Error in parseToken while parsing token")
	}
	if !t.Valid {
		return nil, fmt.Errorf("Error in parseToken, token is no longer valid")
	}
	return t.Claims.(*MyClaims), err
}

//InputValidate checks for invalid input from user
func InputValidate(input string) bool {
	unwantedChar := []string{"\n", "'", "\"", "<", ">", "\t"}
	for _, v := range unwantedChar {
		if strings.Contains(input, v) {
			return false
		}
	}
	return true
}

//EnDecrypt encrypts and decrypts data
func EnDecrypt(msg string) ([]byte, error) {
	encryptionkey = getEnv("ENCRYPTION_KEY")
	encryptKey := []byte(encryptionkey)
	bKey, err := aes.NewCipher(encryptKey[:16])
	if err != nil {
		log.Panic("Error in creating cipher key", err)
	}
	iv := make([]byte, aes.BlockSize)

	s := cipher.NewOFB(bKey, iv)

	buff := &bytes.Buffer{}

	sw := cipher.StreamWriter{
		S: s,
		W: buff,
	}

	_, err = sw.Write([]byte(msg))

	if err != nil {
		return nil, fmt.Errorf("Error in writing to streamWriter")
	}
	return buff.Bytes(), nil
}
