package auth

import (
	datebase "main/database"
	"math/rand"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	_ "github.com/lib/pq"
)

var jwtSecret = []byte(string(rand.Int31()))

type Claims struct {
	User_Id  int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func GenerateToken(username string) (string, error) {

	row := datebase.DB.QueryRow("SELECT id, role FROM dbo.users WHERE name = $1", username)

	var id int
	var role string

	err := row.Scan(&id, &role)
	if err != nil {
		return "", err
	}

	expireTime := time.Now().Add(8 * time.Hour)

	claims := Claims{
		id,
		username,
		role,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "CATSOGRAMM",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {

		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func TokenCheck(header string) (claims *Claims, validity bool, err error) {

	if header == "" {
		return nil, false, nil
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		return nil, false, nil
	}

	if headerParts[0] != "Bearer" {
		return nil, false, nil
	}

	claims, err = ParseToken(headerParts[1])
	if err != nil {
		return nil, false, err
	} else {
		currentTime := time.Now().Unix()
		if claims.ExpiresAt > currentTime {
			return claims, true, nil
		} else {
			return claims, false, nil
		}
	}
}
