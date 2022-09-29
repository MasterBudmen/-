package main

import (
	"log"
	"os"
	_ "strconv"

	database "main/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://catsogramm.web.app")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
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

	//config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://catsogramm.web.app", "https://catsogramm.web.app"}

	//config.AllowAllOrigins = true
	//gr.Use(cors.New(config))

	gr.Use(CORSMiddleware())

	//HEROKU
	r.Run("0.0.0.0:" + os.Getenv("PORT"))
	//LOCAL
	//r.Run(":8080")
}
