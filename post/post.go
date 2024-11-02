package post

import (
	constant "blog/constant"
	db "blog/database"
	ent "blog/ent"

	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/google/uuid"
)

type FrontData struct {
	Title    string    `json:"title"`              // 제목
	Date     string    `json:"date,omitempty"`     // 작성일
	Content  string    `json:"content,omitempty"`  // 마크다운 파일 내용
	Images   []string  `json:"images,omitempty"`   // 이미지 파일 이름 리스트
	TempPath string    `json:"tempPath,omitempty"` // 이미지 임시 저장소 경로
	Category string    `json:"category,omitempty"` // 카테고리
	Tag      []string  `json:"tag,omitempty"`      // 태그 리스트
	ID       uuid.UUID `json:"id,omitempty"`       // ID
}

func SavePost(frontData FrontData) error {

	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return err
	}
	// 현재 시간을 KST 시간대로 변환하여 포맷
	frontData.Date = time.Now().In(loc).Format(time.RFC3339)

	// 프론트엔드에서 받은 게시글 관련 정보를 파일 시스템과 일라스틱서치에 저장하는 함수

	tempPath := constant.MAIN_PATH + "/" + frontData.TempPath
	contentUrl := constant.MAIN_PATH + "/" + frontData.Title
	mdUrl := contentUrl + "/" + frontData.Title + ".md"
	imageUrl := contentUrl + "/images"
	imageUrlForMD := frontData.Title + "/images"

	encodedImageUrlForMD := url.QueryEscape(imageUrlForMD)

	// 1. 본문은 마크다운 파일로 저장
	frontmatter := `---
title: ` + frontData.Title + `
date: ` + frontData.Date + `
tag: ` + strings.Join(frontData.Tag, ",") + `
TocOpen: false,
ShowToc: true,
---
`
	frontData.Content = frontmatter + "\n" + frontData.Content

	// 2. 이미지와 마크다운은 파일시스템에 저장 (constant.MAIN_PATH)
	// title로 폴더 생성
	os.Mkdir(contentUrl, 0755)
	os.Mkdir(imageUrl, 0755)

	// 이미지는 현재 임시 저장소에 저장되어 있으므로, 임시 저장소에서 blog_data로 이동
	for _, image := range frontData.Images {
		os.Rename(tempPath+"/"+image, imageUrl+"/"+image)

		// content 내용도 변경
		frontData.Content = strings.ReplaceAll(frontData.Content, frontData.TempPath+"/"+image, encodedImageUrlForMD+"/"+image)

		// // frontData.Images도 변경
		// frontData.Images[i] = imageUrlForMD + "/" + image
	}

	// 마크다운 파일은 constant.MAIN_PATH에 저장

	os.WriteFile(mdUrl, []byte(frontData.Content), 0644)

	url := "/" + frontData.Title + "/" + frontData.Title + ".md"

	// string to time.Time
	time, err := time.Parse(time.RFC3339, frontData.Date)
	if err != nil {
		return err
	}

	// 3. 마크다운 파일의 메타데이터는 일라스틱서치에 저장
	metadata := ent.Post{
		Title:    frontData.Title,
		Date:     time,
		Content:  url,
		Images:   frontData.Images,
		Category: frontData.Category,
		Tag:      frontData.Tag,
	}

	err = db.SavePost(metadata)
	if err != nil {
		return err
	}

	return nil
}

type Image struct {
	Name   string `json:"name"`
	Source []byte `json:"source"`
}

// 프론트엔드에서 이미지를 byte로 받아서 임시 저장소에 저장하는 함수
func SaveTempImage(images []Image, path string) {
	tempPath := constant.MAIN_PATH + "/temp/" + path
	os.Mkdir(tempPath, 0755)

	// 이미지 저장
	for _, image := range images {
		os.WriteFile(tempPath+"/"+image.Name, image.Source, 0644)
	}
}

func GetPost(row, paginatorNumber int) ([]frontendInfo, error) {

	list, err := db.GetPostList(paginatorNumber, row)
	if err != nil {
		return nil, err
	}

	frontendData, err := readMDFileFromPosts(list)
	if err != nil {
		return nil, err
	}

	return frontendData, nil
}

func DeletePost(title string, postID uuid.UUID) error {

	contentUrl := constant.MAIN_PATH + "/" + title

	// 파일 시스템에서 해당 게시글 삭제
	os.RemoveAll(contentUrl)

	// 일라스틱서치에서 해당 게시글 삭제
	err := db.DeletePost(postID)
	if err != nil {
		return err
	}
	return nil
}

func UpdatePost() {

}

type matter struct {
	Title    string   `yaml:"title"`
	Tags     []string `yaml:"Tags,omitempty"`
	Content  string   `yaml:"content,omitempty"`
	Date     string   `yaml:"date,omitempty"`
	Category string   `yaml:"category,omitempty"`
}

type frontendInfo struct {
	Path   string `json:"path"`
	Matter matter `json:"matter"`
}

func readMDFileFromPosts(list []*ent.Post) ([]frontendInfo, error) {

	var info []frontendInfo

	for _, post := range list {

		title := post.Title

		// 작은따옴표 제거
		title = strings.ReplaceAll(title, "'", "")

		contentUrl := constant.MAIN_PATH + "/" + title + "/" + title + ".md"

		// Read the markdown file
		content, err := os.ReadFile(contentUrl)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		var matter matter
		_, err = frontmatter.Parse(strings.NewReader(string(content)), &matter)
		if err != nil {
			fmt.Println(err)
		}

		// Add the folder information to the list
		info = append(info, frontendInfo{
			Path:   title,
			Matter: matter,
		})
	}
	return info, nil
}

func ReadMDFile(title string) (matter, error) {

	contentUrl := constant.MAIN_PATH + "/" + title + "/" + title + ".md"

	// Read the markdown file
	content, err := os.ReadFile(contentUrl)
	if err != nil {
		fmt.Println(err)
		return matter{}, err
	}

	var matter matter
	_, err = frontmatter.Parse(strings.NewReader(string(content)), &matter)
	if err != nil {
		fmt.Println(err)
	}

	return matter, nil
}

func ConvertLocalFileToDatabaseStruct(matter matter, images []string) ent.Post {

	// string to time.Time
	time, err := time.Parse(time.RFC3339, matter.Date)
	if err != nil {
		fmt.Println(err)
	}

	post := ent.Post{
		Title:    matter.Title,
		Date:     time,
		Content:  matter.Content,
		Images:   images,
		Category: matter.Category,
		Tag:      matter.Tags,
	}

	return post
}
