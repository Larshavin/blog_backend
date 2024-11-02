package handler

import (
	"blog/constant"
	db "blog/database"
	postPack "blog/post"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func TempImageGetHandler() gin.HandlerFunc {
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

func TempImagePostHandler() gin.HandlerFunc {
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

func TempImageDeleteHandler() gin.HandlerFunc {
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

func ListPostHandler() gin.HandlerFunc {
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

		length, err := db.CountAllPosts()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		posts, err := postPack.GetPost(rows, paginatorNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"posts": posts, "length": length})
	}
}

func SavePostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var frontData postPack.FrontData
		err := c.BindJSON(&frontData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		err = postPack.SavePost(frontData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

func DeletePostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var frontData postPack.FrontData
		err := c.BindJSON(&frontData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		err = postPack.DeletePost(frontData.Title, frontData.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

func GetBlogPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := strconv.Atoi(c.Params.ByName("rows"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		title := c.Params.ByName("path")

		var prev, next db.AdjacentPost
		var i int

		prev, next, i, err = db.GetPrevNextPost(title, rows)

		fmt.Println("prev:", prev, "next:", next, "i:", i)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"number": i / rows, "prev": prev, "next": next})
	}
}
