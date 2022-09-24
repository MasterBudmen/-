package main

type User struct {
	Name string `db:"name" json:"name"`
	Role string `db:"role" json:"role"`
}

type UserLogin struct {
	Name     string `form:"name" json:"name" xml:"name"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password"  binding:"required"`
}

type UserRegister struct {
	Name     string `form:"name" json:"name" xml:"name"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password"  binding:"required"`
}

type UserPost struct {
	Post_id string `db:"id" form:"id" json:"id" xml:"id"`
	Name    string `db:"name" form:"name" json:"name" xml:"name"`
	Text    string `db:"text" form:"text" json:"text" xml:"text"  binding:"required"`
}

type UserComment struct {
	Comment_id string `db:"id" form:"id" json:"id" xml:"id"`
	Post_id    string `db:"post_id" form:"post_id" json:"post_id" xml:"post_id"`
	Name       string `db:"name" form:"name" json:"name" xml:"name"`
	Text       string `db:"comment" form:"text" json:"text" xml:"text"  binding:"required"`
}

type UserPostLike struct {
	Name    string `form:"name" json:"name" xml:"name"`
	Post_id string `form:"post_id" json:"post_id" xml:"post_id"  binding:"required"`
}

type UserCommentLike struct {
	Name       string `form:"name" json:"name" xml:"name"`
	Comment_id string `form:"comment_id" json:"comment_id" xml:"post_id"  binding:"required"`
}

type Token struct {
	Token string `form:"token" json:"token" xml:"token"  binding:"required"`
}
