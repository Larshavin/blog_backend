package main

import (
	constant "blog/constant"
	db "blog/database"
	"blog/ent/migrate"
	h "blog/handler"
	postPack "blog/post"
	"context"
	"net"

	"blog/ent"

	_ "github.com/lib/pq"

	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// var (
// 	SSLCRT string = "/etc/letsencrypt/live/3trolls.me/fullchain.pem"
// 	SSLKEY string = "/etc/letsencrypt/live/3trolls.me/privkey.pem"
// )

func loadEnv() {
	envFile := ".env"

	fmt.Println("env : ", os.Getenv("GO_ENV"))

	if os.Getenv("GO_ENV") == "production" {
		envFile = ".env.production"
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file: %v", envFile, err)
	}
}

func main() {

	loadEnv()

	constant.MAIN_PATH = os.Getenv("MAIN_DATA_PATH")
	constant.PORT = os.Getenv("PORT")
	constant.HASH_MIN = os.Getenv("HASH_MIN")
	constant.HASH_MAX = os.Getenv("HASH_MAX")
	constant.MODE = os.Getenv("MODE")

	db.Vars = map[string]string{
		"PG_HOST": os.Getenv("PG_HOST"),
		"PG_PORT": os.Getenv("PG_PORT"),
		"PG_USER": os.Getenv("PG_USER"),
		"PG_DB":   os.Getenv("PG_DB"),
		"PG_PASS": os.Getenv("PG_PASS"),
	}

	fmt.Println("mode : ", constant.MODE)

	// database connection
	postgresDSN := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		db.Vars["PG_HOST"], db.Vars["PG_PORT"], db.Vars["PG_USER"], db.Vars["PG_DB"], db.Vars["PG_PASS"])

	DBClient, err := ent.Open("postgres", postgresDSN)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer DBClient.Close()
	// Run the auto migration tool.
	if err := DBClient.Schema.Create(context.Background(),
		migrate.WithDropColumn(true),
		migrate.WithDropIndex(true),
	); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// 실행 옵션에 --data-save option 을 넣어서 실행하면 로컬에 있는 파일을 데이터베이스에 저장합니다.
	if len(os.Args) > 1 && os.Args[1] == "--data-save" {
		err := CheckLocalPathFileWithDatabase(constant.MAIN_PATH)
		if err != nil {
			log.Fatalf("Error loading %s file: %v", constant.MAIN_PATH, err)
		}
	}

	r := gin.Default()
	r.Use(CORSMiddleware())

	r.Static("/syyang/image", constant.MAIN_PATH).Use(StaticCORSMiddleware())
	r.Static("/syyang/markdown", constant.MAIN_PATH).Use(StaticCORSMiddleware())

	group := r.Group("/syyang")

	group.GET("/tmp/image/:id", h.TempImageGetHandler())
	group.POST("/tmp/image/:id", h.TempImagePostHandler())
	group.DELETE("/tmp/image/:id", h.TempImageDeleteHandler())

	group.GET("/posts/:number/:rows", h.ListPostHandler())
	group.POST("/post", h.SavePostHandler())
	group.DELETE("/post", h.DeletePostHandler())
	group.GET("/post/:path/:rows", h.GetBlogPostHandler())

	group.GET("/tag", h.GetTagsHandler())
	group.GET("/tag/:tag", h.GetTagHandler())

	group.POST("/login", h.PostloginHandler())
	// group.POST("/refresh", user.RefreshTokenHandler(es))
	group.POST("/logout", h.PostLogoutHandler())
	group.POST("/user", h.PostUserHandler())

	group.GET("/check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok", "time": time.Now().Format("2006-01-02 15:04:05")})
	})

	r.Run("0.0.0.0:" + constant.PORT) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	// r.RunTLS(":443", SSLCRT, SSLKEY)
}

func StaticCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

// 172.30.1.0/24 대역을 확인하는 함수
func isAllowedIP(ip string) bool {
	_, ipNet, _ := net.ParseCIDR("172.30.1.0/24")
	parsedIP := net.ParseIP(ip)
	return ipNet.Contains(parsedIP)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowUrlList := []string{"https://blog.3trolls.me", "https://kubesy.com", "http://localhost:5174", "http://172.30.1.40:5174"}
		requestOrigin := c.Request.Header.Get("Origin")
		fmt.Println("requestOrigin:", requestOrigin)
		var allowUrl string

		// IP가 172.30.1.0/24 대역에 있는지 확인 : 개발용
		clientIP := c.ClientIP()
		if requestOrigin == "" {
			if isAllowedIP(clientIP) {
				allowUrl = clientIP
			}
		} else {
			// 허용된 URL 목록에 있는 경우에만 허용
			for _, url := range allowUrlList {
				if requestOrigin == url {
					allowUrl = url
					break
				}
			}
		}

		if allowUrl != "" {
			c.Header("Access-Control-Allow-Origin", allowUrl)
			c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Origin, Accept, X-Requested-With")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		} else {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}

func CheckLocalPathFileWithDatabase(path string) error {

	// check if the file exists in the database
	// path example: ~/blog_data, get list of directory in the path
	blogPostListInLocalPath := []string{}
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	// 디렉토리만 출력합니다.
	for _, file := range files {
		// 폴더가 아니거나 .git 폴더는 제외합니다.
		if file.IsDir() && file.Name() != ".git" && file.Name() != "test" {
			fmt.Println("Directory:", file.Name())
			blogPostListInLocalPath = append(blogPostListInLocalPath, file.Name())
		}
	}

	fmt.Println("blogPostListInLocalPath:", blogPostListInLocalPath)

	for _, blogPostTitle := range blogPostListInLocalPath {
		if !db.CheckPostListInDatabase(blogPostTitle) {
			// save the file to the database
			matter, err := postPack.ReadMDFile(blogPostTitle)
			if err != nil {
				return err
			}
			images, err := findImagesInLocalPost(blogPostTitle)
			if err != nil {
				return err
			}

			fmt.Println("matter:", matter, "images:", images)

			post := postPack.ConvertLocalFileToDatabaseStruct(matter, images)

			fmt.Println("post:", post)

			err = db.SavePost(post)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func findImagesInLocalPost(blogPostTitle string) ([]string, error) {
	// find images in the local post
	images := []string{}
	path := fmt.Sprintf("%s/%s/images", constant.MAIN_PATH, blogPostTitle)
	files, err := os.ReadDir(path)
	if err != nil {
		return images, err
	}
	for _, file := range files {
		if !file.IsDir() {
			images = append(images, file.Name())
		}
	}
	return images, nil
}
