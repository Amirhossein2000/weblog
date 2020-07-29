package main

import (
	"fmt"
	"log"
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
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Println("could not start server err:", err.Error())
	}
}
