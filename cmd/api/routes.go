package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	// grouping to allow seamless switching between version
	v1 := g.Group("/api/v1")
	{
		v1.GET("/ping", app.pingPong)
		v1.POST("/auth/register", app.registerUser)
		v1.POST("/auth/login", app.login)
	}

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	{
		authGroup.GET("/user/me", app.getMyInfo) //get my user info
		authGroup.POST("/user/weight", app.logWeight) //log my weight
		authGroup.GET("/weight", app.getAllWeight) //get all my weight logs
		authGroup.GET("/user", app.getAllUser) //get all user's initial & current weight

	}

	return g
}

func (app *application) pingPong(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
