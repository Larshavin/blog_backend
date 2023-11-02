package main

import (
	"blog/markdown"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	SSLCRT string = "/etc/letsencrypt/live/3trolls.me/fullchain.pem"
	SSLKEY string = "/etc/letsencrypt/live/3trolls.me/privkey.pem"
)

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())

	group := r.Group("/syyang")
	group.Static("/image", "/home/syyang/blog_data")
	group.Static("/markdown", "/home/syyang/blog_data")
	group.GET("/posts/:number/:rows", blogPostsHandler())
	group.GET("/post/:path/:rows", blogPostHandler())

	// r.Run("192.168.15.246:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.RunTLS(":443", SSLCRT, SSLKEY)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		allowUrlList := []string{"http://192.168.15.248:5174", "http://192.168.15.253:5174", "https://kubesy.com", "http://localhost:5174"}
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

		rows, err := strconv.Atoi(c.Params.ByName("rows"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

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
		fmt.Println(length / rows)

		if paginatorNumber-1 > length/rows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Paginator number is too big"})
		} else if paginatorNumber*rows > length {
			c.JSON(http.StatusOK, gin.H{"posts": folders[(paginatorNumber-1)*rows:], "length": length})
		} else {
			c.JSON(http.StatusOK, gin.H{"posts": folders[(paginatorNumber-1)*rows : paginatorNumber*rows], "length": length})
		}
	}
}

type post struct {
	Title           string `json:"title"`
	PaginatorNumber int    `json:"number"`
}

func blogPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := strconv.Atoi(c.Params.ByName("rows"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		path := c.Params.ByName("path")

		folders, err := markdown.FindFolderList("/home/syyang/blog_data")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		length := len(folders)
		for i, folder := range folders {
			if folder.Path == path {
				var prev, next post
				if i == 0 {
					prev.Title = ""
					prev.PaginatorNumber = 0
					next.Title = folders[i+1].Path
					next.PaginatorNumber = (i + 1) / rows
				} else if i == len(folders)-1 {
					prev.Title = folders[i-1].Path
					prev.PaginatorNumber = (i - 1) / rows
					next.Title = ""
					next.PaginatorNumber = length / rows
				} else {
					prev.Title = folders[i-1].Path
					prev.PaginatorNumber = (i - 1) / rows
					next.Title = folders[i+1].Path
					next.PaginatorNumber = (i + 1) / rows
				}
				c.JSON(http.StatusOK, gin.H{"number": i / rows, "prev": prev, "next": next})
				return
			}
		}
	}
}
