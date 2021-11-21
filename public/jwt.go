package public

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

func JwtDecode(tokenStr string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSignKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok {
		return claims, nil
	}

	return nil, errors.New("token is not jwt.StandardClaims")
}

func JwtEncode(claims jwt.StandardClaims) (string, error) {
	singKey := []byte(JwtSignKey)

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return token.SignedString(singKey)
}
