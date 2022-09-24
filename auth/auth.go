package auth

import (
	"database/sql"
	datebase "main/database"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func CheckAuth(username, password string) bool {

	row := datebase.DB.QueryRow("SELECT password FROM dbo.users WHERE name = $1", username)

	var hashPassword string
	err := row.Scan(&hashPassword)
	if err == sql.ErrNoRows {
		return false
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
		if err == nil {
			return true
		} else {
			return false
		}
	}
}
