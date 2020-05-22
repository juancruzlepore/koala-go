package main

import (
	"log"
	"net/http"
	"os"
	_ "github.com/lib/pq"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	//connect()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	dbInstance := connect()
	router := gin.Default()
	////
	// CORS for https://foo.com and https://github.com origins, allowing:
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	corsConfig := cors.Config{
		AllowOrigins:     []string{"https://foranne.herokuapp.com", "http://localhost", "https://localhost"},
		AllowMethods:     []string{"PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true //origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}
	corsInstance := cors.New(corsConfig)
	router.Use(corsInstance)
	router.Use(Options)

	////
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/dates/last", func(c *gin.Context) {
		c.JSON(http.StatusOK, getNextDate(dbInstance))
	})

	router.GET("/movies/all", func(c *gin.Context) {
		c.JSON(http.StatusOK, getMovies(dbInstance))
	})

	router.POST("/movies/add", func(c *gin.Context) {
		var newMovie Movie
        if c.BindJSON(&newMovie) == nil {
			newMovie.CreationDate = time.Now()
			if addMovie(dbInstance, newMovie) {
				c.Status(http.StatusOK)
			} else {
				c.Status(http.StatusConflict)
			}
        } else {
			c.Status(http.StatusBadRequest)
		}		
	})

	router.OPTIONS("/dates/add")

	router.POST("/dates/add", func(c *gin.Context) {
		dateStart := c.PostForm("dateStart")
		dateEnd := c.PostForm("dateEnd")
		if addDate(dbInstance, dateStart, dateEnd) {
			c.Status(http.StatusOK)
		} else {
			c.Status(http.StatusConflict)
		}
	})

	_ = router.Run(":" + port)

}

func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}
}
