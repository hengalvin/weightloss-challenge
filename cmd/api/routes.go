package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	// Optional (silences proxy warning and avoids odd client IP behavior on Render)
	g.SetTrustedProxies(nil)

	// === CORS ===
	g.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://v0-weightloss-challenge-pi.vercel.app",
			"http://localhost:3000", // dev
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // set true ONLY if you send cookies or credentials
		MaxAge:           12 * time.Hour,
	}))

	// (Not strictly necessary with the cors middleware, but harmless)
	g.OPTIONS("/*path", func(c *gin.Context) { c.Status(http.StatusNoContent) })

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
		authGroup.GET("/user/me", app.getMyInfo)      //get my user info
		authGroup.POST("/user/weight", app.logWeight) //log my weight
		authGroup.GET("/weight", app.getAllWeight)    //get all my weight logs
		authGroup.GET("/user", app.getAllUser)        //get all user's initial & current weight

	}

	return g
}

func (app *application) pingPong(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
