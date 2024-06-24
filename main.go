package main

import (
	constant "blog/constant"
	elasticsearch "blog/elastic-search"
	postPack "blog/post"
	user "blog/user"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// var (
// 	SSLCRT string = "/etc/letsencrypt/live/3trolls.me/fullchain.pem"
// 	SSLKEY string = "/etc/letsencrypt/live/3trolls.me/privkey.pem"
// )

func main() {

	// get parameter from .env file
	var myEnv map[string]string
	myEnv, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	constant.MAIN_PATH = myEnv["MAIN_DATA_PATH"]
	constant.PORT = myEnv["PORT"]
	constant.HASH_MIN = myEnv["HASH_MIN"]
	constant.HASH_MAX = myEnv["HASH_MAX"]
	constant.MODE = myEnv["MODE"]

	if myEnv["ELASTICSEARCH_TLS"] != "" {
		constant.ELASTICSEARCH_TLS = myEnv["ELASTICSEARCH_TLS"]
	} else if os.Getenv("ELASTICSEARCH_TLS") != "" {
		constant.ELASTICSEARCH_TLS = os.Getenv("ELASTICSEARCH_TLS")
	} else {
		// program exit
		log.Fatal("ELASTICSEARCH_TLS is not set")
	}

	es := elasticsearch.Client{}
	es.Connect(constant.MODE)

	r := gin.Default()
	r.Use(CORSMiddleware())

	group := r.Group("/syyang")

	group.Static("/image", constant.MAIN_PATH)
	group.Static("/markdown", constant.MAIN_PATH)

	group.GET("/tmp/image/:id", tmpImageGetHandler())
	group.POST("/tmp/image/:id", tmpImagePostHandler())
	group.DELETE("/tmp/image/:id", tmpImageDeleteHandler())

	group.GET("/posts/:number/:rows", listPostHandler(es))
	group.POST("/post", savePostHandler(es))
	group.DELETE("/post", deletePostHandler(es))
	group.GET("/post/:path/:rows", getBlogPostHandler(es))

	group.POST("/login", postloginHandler(es))
	// group.POST("/refresh", user.RefreshTokenHandler(es))
	group.POST("/logout", postLogoutHandler(es))

	group.GET("/check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok", "time": time.Now().Format("2006-01-02 15:04:05")})
	})

	r.Run("0.0.0.0:" + constant.PORT) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	// r.RunTLS(":443", SSLCRT, SSLKEY)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		allowUrlList := []string{"https://blog.3trolls.me", "https://kubesy.com", "http://localhost:5174", "http://172.30.1.40:5174"}
		var allowUrl string
		for _, url := range allowUrlList {
			if c.Request.Header.Get("Origin") == url {
				allowUrl = url
				break
			}
		}

		c.Header("Access-Control-Allow-Origin", allowUrl)
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Origin, Accept, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func tmpImageGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")
		tmpPath := fmt.Sprintf("%s/%s", constant.MAIN_PATH, id)
		err := os.MkdirAll(tmpPath, 0755)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"path": tmpPath})
	}
}

func tmpImagePostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
			return
		}

		id := c.Params.ByName("id")
		tmpPath := fmt.Sprintf("%s/%s", constant.MAIN_PATH, id)
		fileName := fmt.Sprintf("%s/%s", tmpPath, file.Filename)
		if err := c.SaveUploadedFile(file, fileName); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"message": "ok", "path": id + "/" + file.Filename, "size": file.Size, "name": file.Filename})
	}
}

func tmpImageDeleteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")
		tmpPath := fmt.Sprintf("%s/%s", constant.MAIN_PATH, id)
		err := os.RemoveAll(tmpPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

/* */

func listPostHandler(es elasticsearch.Client) gin.HandlerFunc {
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

		length, err := es.CountAllPosts()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		posts, err := postPack.GetPost(rows, paginatorNumber, &es)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"posts": posts, "length": length})
	}
}

func savePostHandler(es elasticsearch.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var frontData postPack.FrontData
		err := c.BindJSON(&frontData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		err = postPack.SavePost(frontData, &es)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

func deletePostHandler(es elasticsearch.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var frontData postPack.FrontData
		err := c.BindJSON(&frontData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		err = postPack.DeletePost(frontData.Title, frontData.ID, &es)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

func getBlogPostHandler(es elasticsearch.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := strconv.Atoi(c.Params.ByName("rows"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		title := c.Params.ByName("path")

		var prev, next elasticsearch.AdjacentPost
		var i int

		prev, next, i, err = es.GetPrevNextPost(title, rows)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"number": i / rows, "prev": prev, "next": next})
	}
}

func postloginHandler(es elasticsearch.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		var loginData elasticsearch.User
		err := c.BindJSON(&loginData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		userInfo, err := user.CheckPassword(loginData.Email, loginData.Password, es)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		accessToken, refreshToken, err := user.GenerateToken(userInfo, es)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		// marshal data
		data := gin.H{"accessToken": accessToken, "refreshToken": refreshToken}

		c.JSON(http.StatusOK, data)
	}
}

func postLogoutHandler(es elasticsearch.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		type dataType struct {
			AccessToken string `json:"accessToken"`
		}
		var json dataType

		// decode token and get email
		err := c.BindJSON(&json)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		data, err := user.DecodeToken(json.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		userEmail := data.Email

		// delete two tokens from elastic search

		err = es.DeleteToken(userEmail)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"message": "logout success"})
	}
}
