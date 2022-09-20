package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	_ "strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	c_host     string
	c_port     string
	c_user     string
	c_password string
	c_dbname   string
	c_APIKey   string

	db *sql.DB
)

type User struct {
	Name     string `db:"name" json:"name"`
	Password string `db:"password" json:"password"`
	Role     string `db:"role" json:"role"`
}

type APIKey struct {
	Key string `form:"APIKey" json:"APIKey" xml:"APIKey"  binding:"required"`
}

type UserLogin struct {
	Name     string `form:"name" json:"name" xml:"name"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password"  binding:"required"`
	Key      string `form:"APIKey" json:"APIKey" xml:"APIKey"  binding:"required"`
}

type UserRegister struct {
	Name     string `form:"name" json:"name" xml:"name"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password"  binding:"required"`
	Role     string `form:"role" json:"role" xml:"role"  binding:"required"`
	Key      string `form:"APIKey" json:"APIKey" xml:"APIKey"  binding:"required"`
}

func initDb() *sql.DB {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func GetUsers(c *gin.Context) {

	var Key APIKey

	if err := c.ShouldBind(&Key); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if Key.Key != c_APIKey {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	var users []User
	rows, err := db.Query("SELECT name, password, role FROM dbo.users")
	if err != nil {
		c.JSON(404, gin.H{"error": "unknown error"})
	}

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Name, &user.Password, &user.Role); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": users,
	})
}

func Login(c *gin.Context) {

	var user UserLogin

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Key != c_APIKey {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	row := db.QueryRow("SELECT role FROM dbo.users WHERE name = $1 AND password = $2", user.Name, user.Password)

	var role string
	err := row.Scan(&role)
	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "user is not registered"})
	} else {
		c.JSON(http.StatusOK, gin.H{"role": role})
	}

	row.Scan()
}

func Register(c *gin.Context) {

	var user UserRegister

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Key != c_APIKey {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	row := db.QueryRow("SELECT role FROM dbo.users WHERE name = $1", user.Name)

	var role string
	err := row.Scan(&role)
	if err != sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "user is already registered"})
	} else {
		db.Exec("INSERT INTO dbo.users (name, password, role) VALUES ($1, $2, $3)", user.Name, user.Password, user.Role)
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK})
	}

	row.Scan()
}

func main() {

	err := godotenv.Load("consts.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	c_host = os.Getenv("host")
	c_port = os.Getenv("port")
	c_user = os.Getenv("user")
	c_password = os.Getenv("password")
	c_dbname = os.Getenv("dbname")
	c_APIKey = os.Getenv("APIKey")

	db = initDb()

	r := gin.Default()

	gr := r.Group("/api")

	gr.POST("/users", GetUsers)
	gr.POST("/login", Login)
	gr.POST("/register", Register)

	r.Run("0.0.0.0:" + os.Getenv("PORT"))
}
