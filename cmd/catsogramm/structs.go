package main

//import (
//	"database/sql"
//)

type User struct {
	Id    string `db:"id" json:"id"`
	Login string `db:"login" json:"login"`
	Role  string `db:"role" json:"role"`
}

type UserLogin struct {
	Login    string `form:"login" json:"login" xml:"login"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password"  binding:"required"`
}

type UserRegister struct {
	Login    string `form:"login" json:"login" xml:"login"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password"  binding:"required"`
}

type UserPostGET struct {
	Post_id string `db:"id" form:"id" json:"id" xml:"id"`
	User_Id string `db:"user_Id" form:"user_Id" json:"user_Id" xml:"user_Id"`
	Login   string `db:"login" form:"login" json:"login" xml:"login"`
	Text    string `db:"text" form:"text" json:"text" xml:"text"  binding:"required"`
	Image   string `db:"Image" form:"Image" json:"Image" xml:"Image"`
}

type UserPostPOST struct {
	Text     string `db:"text" form:"text" json:"text" xml:"text"  binding:"required"`
	Image_id string `db:"image_id" form:"image_id" json:"image_id" xml:"image_id"`
}

type Image struct {
	Image_id string `db:"id" form:"id" json:"id" xml:"id"`
	Image    string `db:"image" form:"image" json:"image" xml:"image"`
}

type UserComment struct {
	Comment_id string `db:"id" form:"id" json:"id" xml:"id"`
	Post_id    string `db:"post_id" form:"post_id" json:"post_id" xml:"post_id"`
	Login      string `db:"login" form:"login" json:"login" xml:"login"`
	Text       string `db:"comment" form:"text" json:"text" xml:"text"  binding:"required"`
}

type UserPostLike struct {
	Login   string `form:"login" json:"login" xml:"login"`
	Post_id string `form:"post_id" json:"post_id" xml:"post_id"  binding:"required"`
}

type UserCommentLike struct {
	Login      string `form:"login" json:"login" xml:"login"`
	Comment_id string `form:"comment_id" json:"comment_id" xml:"post_id"  binding:"required"`
}

type Token struct {
	Token string `form:"token" json:"token" xml:"token"  binding:"required"`
}
