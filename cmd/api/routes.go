package main

import (
	"net/http"
	"net/url"
	"regexp"
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
		AllowOriginFunc: func(origin string) bool {
			// Parse the URL to avoid false positives
			u, err := url.Parse(origin)
			if err != nil {
				return false
			}

			// Regex to match the preview URLs you described
			matched, _ := regexp.MatchString(
				`^preview-weight-loss-app-[a-zA-Z0-9]+\.vusercontent\.net$`,
				u.Host,
			)
			if matched {
				return true
			}

			// Explicitly allow production site too
			if u.Host == "v0-weightloss-challenge-pi.vercel.app" {
				return true
			}

			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
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
