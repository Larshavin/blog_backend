package database

import (
	ent "blog/ent"
	"blog/ent/post"
	"blog/ent/tag"
	"context"
)

func GetPostListByTag(t string) ([]*ent.Post, error) {

	client := connection()
	defer client.Close()
	ctx := context.Background()

	posts, err := client.Post.Query().Where(post.HasTagsWith(tag.Name(t))).All(ctx)

	if err != nil {
		return nil, err
	}

	return posts, nil
}
