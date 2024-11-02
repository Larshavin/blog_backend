package handler

import (
	db "blog/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTagsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tags, err := db.GetTags()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"tags": tags})
	}
}

func GetTagHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tag := c.Param("tag")

		// tag 에 연결된 post 목록 가져오기
		postList, err := db.GetPostListByTag(tag)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"tag": tag, "postList": postList})
	}
}
