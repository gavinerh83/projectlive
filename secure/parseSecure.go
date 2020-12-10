//Package secure contains functions that provide security to overall application.
//It contains functions that generate and validate JWT for each session login.
package secure

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

//MyClaims is the claims struct that is used in the generation of JWT
type MyClaims struct {
	jwt.StandardClaims
	SessionID string
}

var (
	key = []byte("this is JWT key")
)

//GenerateJWT creates a JWT for each user session login
func GenerateJWT(c *MyClaims) (string, error) {
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
