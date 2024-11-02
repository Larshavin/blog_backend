package database

import (
	ent "blog/ent"
	"blog/ent/post"
	"blog/ent/tag"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type AdjacentPost struct {
	Title           string `json:"title"`
	PaginatorNumber int    `json:"number"`
}

func CountAllPosts() (int, error) {
	client := connection()
	defer client.Close()
	ctx := context.Background()

	// return the number of posts
	n, err := client.Post.Query().Count(ctx)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func GetPrevNextPost(title string, rows int) (AdjacentPost, AdjacentPost, int, error) {

	client := connection()
	defer client.Close()
	ctx := context.Background()

	// get the list of posts
	posts, err := client.Post.Query().Order(ent.Desc(post.FieldDate)).All(ctx)
	if err != nil {
		return AdjacentPost{}, AdjacentPost{}, 0, err
	}

	// find the index of the post
	var i int
	for i = 0; i < len(posts); i++ {
		if posts[i].Title == title {
			break
		}
	}

	// find the previous and next post
	if i > 0 {
		prev := AdjacentPost{
			Title:           posts[i-1].Title,
			PaginatorNumber: (i - 1) / rows,
		}
		if i < len(posts)-1 {
			next := AdjacentPost{
				Title:           posts[i+1].Title,
				PaginatorNumber: (i + 1) / rows,
			}
			return prev, next, i, nil
		}
		return prev, AdjacentPost{}, i, nil
	} else if i < len(posts)-1 {
		next := AdjacentPost{
			Title:           posts[i+1].Title,
			PaginatorNumber: (i + 1) / rows,
		}
		return AdjacentPost{}, next, i, nil
	}

	return AdjacentPost{}, AdjacentPost{}, 0, nil
}

func checkTag(tags []string) []*ent.Tag {

	client := connection()
	defer client.Close()
	ctx := context.Background()

	// 우선 태그를 처리합니다.
	var tagEntities []*ent.Tag
	for _, tagName := range tags {
		// 태그가 이미 존재하는지 확인합니다.
		tag, err := client.Tag.
			Query().
			Where(tag.Name(tagName)).
			Only(ctx)

		if ent.IsNotFound(err) {
			// 태그가 없으면 새로 생성합니다.
			tag, err = client.Tag.
				Create().
				SetName(tagName).
				Save(ctx)
			if err != nil {
				return nil
			}
		} else if err != nil {
			return nil
		}

		// 태그 엔티티를 리스트에 추가합니다.
		tagEntities = append(tagEntities, tag)
	}
	return tagEntities
}

func SavePost(post ent.Post) error {
	client := connection()
	defer client.Close()
	ctx := context.Background()

	tags := checkTag(post.Tag)

	_, err := client.Post.Create().
		SetTitle(post.Title).
		SetContent(post.Content).
		SetDate(post.Date).
		SetCategory(post.Category).
		AddTags(tags...).
		SetImages(post.Images).
		Save(ctx)
	if err != nil {
		fmt.Println("게시글 저장 실패")
		return err
	}

	return nil
}

func GetPostList(paginatorNumber int, rows int) ([]*ent.Post, error) {

	client := connection()
	defer client.Close()
	ctx := context.Background()

	// get the list of posts with pagination information
	list, err := client.Post.Query().
		Limit(rows).
		Offset((paginatorNumber - 1) * rows).
		Order(ent.Desc(post.FieldDate)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func DeletePost(id uuid.UUID) error {

	client := connection()
	defer client.Close()
	ctx := context.Background()

	_, err := client.Post.Delete().Where(post.IDEQ(id)).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func GetPost(id uuid.UUID) (*ent.Post, error) {

	client := connection()
	defer client.Close()
	ctx := context.Background()

	post, err := client.Post.Query().Where(post.IDEQ(id)).Only(ctx)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func UpdatePost(p ent.Post) error {

	client := connection()
	defer client.Close()
	ctx := context.Background()

	tags := checkTag(p.Tag)

	_, err := client.Post.Update().
		Where(post.IDEQ(p.ID)).
		SetTitle(p.Title).
		SetContent(p.Content).
		SetDate(p.Date).
		SetCategory(p.Category).
		SetImages(p.Images).
		AddTags(tags...).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func GetPostByCategory(category string) ([]*ent.Post, error) {

	client := connection()
	defer client.Close()
	ctx := context.Background()

	posts, err := client.Post.Query().Where(post.CategoryEQ(category)).All(ctx)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func GetPostByTag(tagName string) ([]*ent.Post, error) {

	client := connection()
	defer client.Close()
	ctx := context.Background()

	posts, err := client.Post.Query().Where(post.HasTagsWith(tag.Name(tagName))).All(ctx)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func GetPostByDate(date string) ([]*ent.Post, error) {
	return nil, nil
}

func GetPostBySearch(search string) ([]*ent.Post, error) {

	return nil, nil
}

func GetTags() ([]*ent.Tag, error) {

	client := connection()
	defer client.Close()
	ctx := context.Background()

	tags, err := client.Tag.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func CheckPostListInDatabase(title string) bool {
	client := connection()
	defer client.Close()
	ctx := context.Background()

	_, err := client.Post.Query().Where(post.TitleEQ(title)).Only(ctx)
	return err == nil // if err is nil, the post exists in the database
}
