package main

import (
	"log"
	_ "strconv"

	"main/internal/database"
	"main/internal/restapi"

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
	defer database.DB.Close()

	r := gin.Default()

	gr := r.Group("/api")

	r.Use(CORSMiddleware())
	gr.Use(CORSMiddleware())

	gr.GET("/users", restapi.GetUsers)
	gr.POST("/users/register", restapi.Register)
	gr.POST("/users/login", restapi.Login)

	gr.POST("/images", restapi.UploadImage)
	gr.GET("/images/:id", restapi.DownloadImage)

	//gr.POST("/tokencheck", ReadToken)

	gr.GET("/comments", restapi.GetComments)
	gr.GET("/posts", restapi.GetPosts)

	gr.POST("/posts", restapi.Post)
	gr.POST("/comments", restapi.Comment)

	gr.POST("/posts/:id/like", restapi.Like_Post)
	gr.POST("/comments/:id/like", restapi.Like_Comment)

	gr.GET("/app-check", restapi.AppCheck)
	gr.GET("/db-check", restapi.DBCheck)

	gr.StaticFile("./doc.json", "../../api/swagger.json")
	r.Static("/swaggerui/", "../../web/swaggerui")

	//HEROKU
	//r.Run("0.0.0.0:" + os.Getenv("PORT"))
	//LOCAL
	r.Run(":8080")
}
