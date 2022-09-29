package main

import (
	"log"
	"os"
	_ "strconv"

	database "main/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func GetUserIdByName(user User) int {
	return 0
}

func main() {

	err := godotenv.Load("consts.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.DB = database.InitDb()

	r := gin.Default()

	gr := r.Group("/api")

	gr.GET("/users", GetUsers)
	gr.POST("/users/register", Register)
	gr.POST("/users/login", Login)

	//gr.POST("/tokencheck", TokenParse)

	gr.GET("/comments", GetComments)
	gr.GET("/posts", GetPosts)

	gr.POST("/posts/create", Post)
	gr.POST("/comments/create", Comment)

	gr.POST("/posts/like", Like_Post)
	gr.POST("/comments/like", Like_Comment)

	gr.GET("/app-check", AppCheck)
	gr.GET("/db-check", DBCheck)

	r.Static("/swaggerui/", "swaggerui")

	gr.Use(cors.Default())
	//HEROKU
	r.Run("0.0.0.0:" + os.Getenv("PORT"))
	//LOCAL
	//r.Run(":8080")
}
