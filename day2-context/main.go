package main

import (
	"Gee/day2-context/gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/", func(context *gee.Context) {
		context.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	r.GET("/hello", func(context *gee.Context) {
		context.String(http.StatusOK, "hello %s,you're at %s\n", context.Query("name"), context.Path)
	})
	r.POST("/login", func(context *gee.Context) {
		context.JSON(http.StatusOK, gee.H{
			"username": context.PostForm("username"),
			"password": context.PostForm("password"),
		})
	})
	r.Run(":9999")
}
