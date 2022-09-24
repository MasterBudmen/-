package main

import (
	"log"
	"os"
	_ "strconv"

	database "main/database"

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

	gr.POST("/register", Register)
	gr.POST("/login", Login)
	gr.POST("/tokencheck", TokenParse)

	gr.GET("/users", GetUsers)
	gr.GET("/comments", GetComments)
	gr.GET("/posts", GetPosts)

	gr.POST("/post", Post)
	gr.POST("/comment", Comment)

	gr.POST("/post_like", Like_Post)
	gr.POST("/comment_like", Like_Comment)

	gr.GET("/app-check", AppCheck)
	gr.GET("/db-check", DBCheck)

	r.Static("/swaggerui/", "swaggerui")

	//HEROKU
	r.Run("0.0.0.0:" + os.Getenv("PORT"))
	//LOCAL
	//r.Run(":8080")
}
