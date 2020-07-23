package main

import (
	"fmt"
	"net/http"
	"weblog/controllers/article"
	"weblog/controllers/auth"
	"weblog/controllers/comment"
	"weblog/controllers/timeline"
	"weblog/controllers/user"
)

func main() {
	http.HandleFunc("/login", auth.LoginController)
	http.HandleFunc("/register", auth.Register)

	http.HandleFunc("/user", user.UserController)
	http.HandleFunc("/article", article.ArticleController)
	http.HandleFunc("/comment", comment.CommentController)

	http.HandleFunc("/timeline", timeline.TimelineController)

	port := ":8080"
	fmt.Println("server started on port", port)
	http.ListenAndServe(port, nil)
}
