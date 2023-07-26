package main

import (
	"blog/markdown"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.Static("/image", "/home/syyang/blog_data")
	r.Static("/markdown", "/home/syyang/blog_data")
	r.GET("/posts/:number", blogPostsHandler())

	r.Run("192.168.15.246:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		allowUrlList := []string{"http://192.168.15.248:5174", "https://kubesy.com", "http://localhost:5174"}
		var allowUrl string
		for _, url := range allowUrlList {
			if c.Request.Header.Get("Origin") == url {
				allowUrl = url
				break
			}
		}

		// c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
		// c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", allowUrl)
		c.Header("Access-Control-Allow-Methods", "GET, DELETE, POST")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func blogPostsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		number := c.Params.ByName("number")
		paginatorNumber, err := strconv.Atoi(number)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		folders, err := markdown.FindFolderList("/home/syyang/blog_data")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		length := len(folders)
		fmt.Println(length / 10)

		if paginatorNumber-1 > length/10 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Paginator number is too big"})
		} else if paginatorNumber*10 > length {
			c.JSON(http.StatusOK, folders[(paginatorNumber-1)*10:])
		} else {
			c.JSON(http.StatusOK, folders[(paginatorNumber-1)*10:paginatorNumber*10])
		}
	}
}
