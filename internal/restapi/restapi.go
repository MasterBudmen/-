package restapi

import (
	"database/sql"
	"main/internal/auth"
	"main/internal/model"
	"net/http"
	_ "strconv"
	"time"

	database "main/internal/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user model.UserLoginRequest

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	row := database.DB.QueryRow("SELECT role FROM dbo.users WHERE name = $1", user.Login)

	passwordhash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var role string
	err = row.Scan(&role)
	if err != sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User is already registered"})
	} else {
		database.DB.Exec("INSERT INTO dbo.users (name, password, role, registered_at) VALUES ($1, $2, $3, $4)", user.Login, string(passwordhash), "user", time.Now())
		c.JSON(http.StatusOK, gin.H{"message": "User has been registered"})
	}
}

func AppCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Good"})
}

func DBCheck(c *gin.Context) {
	err := database.DB.Ping()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Good"})
	}
}

func Post(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	claims, validity, _ := auth.TokenCheck(authHeader)

	if validity {
		var post model.UserPostRequest

		if err := c.ShouldBind(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		var err error
		if post.Image_id == "" {
			_, err = database.DB.Exec("INSERT INTO dbo.posts (user_id, txt, created_at, updated_at) VALUES ($1, $2, $3, $3)", claims.User_Id, post.Text, time.Now())
		} else {
			_, err = database.DB.Exec("INSERT INTO dbo.posts (user_id, txt, created_at, updated_at, image_id) VALUES ($1, $2, $3, $3, $4)", claims.User_Id, post.Text, time.Now(), post.Image_id)
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "post has been created"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
	}
}

func GetPosts(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")

	_, validity, _ := auth.TokenCheck(authHeader)

	if validity {
		var posts []model.UserPostResponse
		limit, exists := c.GetQuery("limit")
		if !exists {
			limit = "10"
		}

		offset, exists := c.GetQuery("offset")

		if !exists {
			offset = "0"
		}

		var rows *sql.Rows
		var err error

		user_id, exists := c.GetQuery("user_id")
		if !exists {
			rows, err = database.DB.Query("SELECT p.id, p.user_id, s.name, p.txt as text, i.image FROM dbo.posts p JOIN dbo.users s ON p.user_id = s.id LEFT JOIN dbo.images i ON p.image_id = i.id ORDER BY p.created_at DESC LIMIT $1 OFFSET $2", limit, offset)
		} else {
			rows, err = database.DB.Query("SELECT p.id, p.user_id, s.name, p.txt as text, i.image FROM dbo.posts p JOIN dbo.users s ON p.user_id = s.id LEFT JOIN dbo.images i ON p.image_id = i.id WHERE p.user_id = $1 ORDER BY p.created_at DESC LIMIT $2 OFFSET $3", user_id, limit, offset)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}

		var img sql.NullString

		for rows.Next() {
			var post model.UserPostResponse
			if err := rows.Scan(&post.Post_id, &post.User_Id, &post.Login, &post.Text, &img); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			if img.Valid {
				post.Image = img.String
			}
			posts = append(posts, post)
		}
		c.JSON(http.StatusOK, gin.H{
			"posts": posts,
		})

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
	}
}

func Comment(c *gin.Context) {
	var comment model.UserComment

	authHeader := c.GetHeader("Authorization")

	if err := c.ShouldBind(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	claims, validity, _ := auth.TokenCheck(authHeader)

	if validity {
		_, err := database.DB.Exec("INSERT INTO dbo.comments (post_id, user_id, comment, created_at, updated_at) VALUES ($1, $2, $3, $4, $4)", comment.Post_id, claims.User_Id, comment.Text, time.Now())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "comment has been created"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
	}
}

func GetComments(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	_, validity, _ := auth.TokenCheck(authHeader)

	if validity {
		var comments []model.UserComment

		post_id, exists := c.GetQuery("post_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "post_id is not specified"})
			return
		}

		limit, exists := c.GetQuery("limit")
		if !exists {
			limit = "10"
		}

		offset, exists := c.GetQuery("offset")

		if !exists {
			offset = "0"
		}

		rows, err := database.DB.Query("SELECT c.id, c.post_id, s.name, c.comment as text FROM dbo.comments c JOIN dbo.users s ON c.user_id = s.id WHERE c.post_id = $1 ORDER BY c.created_at DESC LIMIT $2 OFFSET $3", post_id, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}

		for rows.Next() {
			var comment model.UserComment
			if err := rows.Scan(&comment.Comment_id, &comment.Post_id, &comment.Login, &comment.Text); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			comments = append(comments, comment)
		}
		c.JSON(http.StatusOK, gin.H{
			"comments": comments,
		})

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
	}
}

func UploadImage(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	claims, validity, _ := auth.TokenCheck(authHeader)

	if validity {
		var img model.Image

		if err := c.ShouldBind(&img); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		stmt, err := database.DB.Prepare("INSERT INTO dbo.images (image, created_at, user_id) VALUES ($1, $2, $3) RETURNING id")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		defer stmt.Close()

		var id int

		err = stmt.QueryRow(
			img.Image,
			time.Now(),
			claims.User_Id,
		).Scan(&id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message":  "image has been uploaded",
				"image_id": id})
		}

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
	}
}

func DownloadImage(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	_, validity, _ := auth.TokenCheck(authHeader)

	if validity {

		id := c.Param("id")
		var img string

		row := database.DB.QueryRow("SELECT image FROM dbo.images WHERE id = $1", id)

		err := row.Scan(&img)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "image not found"})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"image": img})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
	}
}

func Like_Post(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	claims, validity, _ := auth.TokenCheck(authHeader)

	if validity {
		id := c.Param("id")
		row := database.DB.QueryRow("SELECT user_id FROM dbo.posts_likes WHERE post_id = $1 AND user_id = $2", id, claims.User_Id)

		var user_id string
		err := row.Scan(&user_id)
		if err != sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"message": "there is like already"})
		} else {
			_, err := database.DB.Exec("INSERT INTO dbo.posts_likes (user_id, post_id, created_at) VALUES ($1, $2, $3)", claims.User_Id, id, time.Now())
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "like has been created"})
			}
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
	}
}

func Like_Comment(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")

	claims, validity, _ := auth.TokenCheck(authHeader)

	if validity {

		id := c.Param("id")

		row := database.DB.QueryRow("SELECT user_id FROM dbo.comments_likes WHERE comment_id = $1 AND user_id = $2", id, claims.User_Id)

		var user_id string
		err := row.Scan(&user_id)
		if err != sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"message": "there is like already"})
		} else {
			_, err := database.DB.Exec("INSERT INTO dbo.comments_likes (user_id, comment_id, created_at) VALUES ($1, $2, $3)", claims.User_Id, id, time.Now())
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "like has been created"})
			}
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
	}
}

func GetUsers(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")

	_, validity, _ := auth.TokenCheck(authHeader)

	if validity {
		var users []model.User
		limit, exists := c.GetQuery("limit")
		if !exists {
			limit = "10"
		}

		offset, exists := c.GetQuery("offset")

		if !exists {
			offset = "0"
		}

		rows, err := database.DB.Query("SELECT id, name, role FROM dbo.users ORDER BY registered_at DESC LIMIT $1 OFFSET $2", limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}

		for rows.Next() {
			var user model.User
			if err := rows.Scan(&user.Id, &user.Login, &user.Role); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			users = append(users, user)
		}
		c.JSON(http.StatusOK, gin.H{
			"users": users,
		})

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
	}
}

func Login(c *gin.Context) {

	var user model.UserLoginRequest

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if auth.CheckAuth(user.Login, user.Password) {
		token, err := auth.GenerateToken(user.Login)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Error generating token"})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "Authorized",
				"token":   token,
			})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "user is not registered or invalid password",
		})
	}
}

func ReadToken(c *gin.Context) {
	var token model.Token
	if err := c.ShouldBind(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	} else {
		claims, validity, err := auth.TokenCheck(token.Token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			code := http.StatusUnauthorized
			if validity {
				code = http.StatusOK
			}

			c.JSON(code, gin.H{
				"validity":  validity,
				"user_id":   claims.User_Id,
				"username":  claims.Username,
				"ExpiredAt": claims.ExpiresAt,
				"Now":       time.Now().Unix(),
			})
		}
	}
}
