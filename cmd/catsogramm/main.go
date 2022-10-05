package main

import (
	"log"
	_ "strconv"

	database "main/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "https://catsogramm.web.app")
		c.Header("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {

	err := godotenv.Load("./../../configs/consts.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.DB = database.InitDb()

	r := gin.Default()

	gr := r.Group("/api")

	r.Use(CORSMiddleware())
	gr.Use(CORSMiddleware())

	gr.GET("/users", GetUsers)
	gr.POST("/users/register", Register)
	gr.POST("/users/login", Login)

	gr.POST("/images", UploadImage)
	gr.GET("/images/:id", DownloadImage)

	//gr.POST("/tokencheck", TokenParse)

	gr.GET("/comments", GetComments)
	gr.GET("/posts", GetPosts)

	gr.POST("/posts", Post)
	gr.POST("/comments", Comment)

	gr.POST("/posts/:id/like", Like_Post)
	gr.POST("/comments/:id/like", Like_Comment)

	gr.GET("/app-check", AppCheck)
	gr.GET("/db-check", DBCheck)

	gr.StaticFile("./doc.json", "../../api/swagger.json")
	r.Static("/swaggerui/", "../../web/swaggerui")

	//HEROKU
	//r.Run("0.0.0.0:" + os.Getenv("PORT"))
	//LOCAL
	r.Run(":8080")
}
