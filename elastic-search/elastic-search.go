package elasticsearch

// This is a package for elasticsearch client

import (
	"blog/constant"
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

type Client struct {
	Client *elasticsearch.Client
}

type Hit struct {
	Source interface{} `json:"_source"`
}

type Hits struct {
	Total struct {
		Value int `json:"value"`
	} `json:"total"`
	Hits []Hit `json:"hits"`
}

type SearchResult struct {
	Hits Hits `json:"hits"`
}

func (c *Client) Connect(mod string) {

	// // Load the TLS certificate from file
	// cert, err := os.ReadFile("elastic-search/elastic_tls.crt")
	// if err != nil {
	// 	log.Fatalf("Failed to read certificate: %s", err)
	// }

	// Load the TLS certificate from environment variable
	cert := []byte(constant.ELASTICSEARCH_TLS)

	// Create a new cert pool and add our cert to it
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		log.Fatal("Failed to append cert")
	}

	// Create a TLS Config
	tlsConfig := &tls.Config{
		RootCAs: certPool,
		// InsecureSkipVerify: true, // use this line if you want to ignore SSL verification (not recommended for production)
	}

	// Set up the HTTP transport and include the TLS config
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	var url string
	if mod == constant.DEV {
		url = "https://172.30.1.154:9200/"
	} else {
		url = "https://quickstart-es-http.blog.svc:9200/"
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			url,
		},
		APIKey:    "Vk5NNlhvOEIxZkM3M19vdVR0bzE6WXdxdHgzYXZSNUNLc2pOOWVNWUJ6UQ==",
		Transport: transport,
	}

	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	c.Client = esClient
}

func (c *Client) ClientTest() {
	infores, err := c.Client.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	fmt.Println(infores)
}

// 데이터 모델 정의
type Post struct {
	Title    string   `json:"title"`              // 제목
	Date     string   `json:"date,omitempty"`     // 작성일
	Content  string   `json:"content,omitempty"`  // 마크다운 파일 경로
	Images   []string `json:"images,omitempty"`   // 이미지 파일 경로 리스트
	Category string   `json:"category,omitempty"` // 카테고리
	Tag      []string `json:"tag,omitempty"`      // 태그 리스트
}

func (c *Client) SavePost(post Post) error {

	data, _ := json.Marshal(post)
	response, err := c.Client.Index("kubesy", bytes.NewReader(data))
	if err != nil {
		return err
	}

	// Check response status
	if response.IsError() {
		err := fmt.Errorf("error: %s", response.String())
		return err
	}

	return nil
}

func (c *Client) GetPostList(start, count int) ([]interface{}, error) {

	// 검색 쿼리 실행
	query := fmt.Sprintf(`{
		"from": %d,
		"size": %d,
		"sort": [
			{ "date": { "order": "desc" } }
		],
		"query": {
			"match_all": {}
		}
	}`, start, count)

	searchResult, err := c.Client.Search(
		c.Client.Search.WithIndex("kubesy"),
		c.Client.Search.WithBody(strings.NewReader(query)),
	) // 쿼리 실행
	if err != nil {
		return nil, err
	}

	if searchResult.IsError() {
		err := fmt.Errorf("error: %s", searchResult.String())
		return nil, err
	}

	// decode the response
	var response map[string]interface{}

	if err := json.NewDecoder(searchResult.Body).Decode(&response); err != nil {
		return nil, err
	}

	// Get only hits slice
	hits := response["hits"].(map[string]interface{})["hits"].([]interface{})

	return hits, nil
}

func (c *Client) DeletePost(id string) error {
	response, err := c.Client.Delete("kubesy", id)
	if err != nil {
		return err
	}
	fmt.Println(response)
	return nil
}

func (c *Client) GetPostIDByTitle(title string) (string, error) {

	// 검색 쿼리 실행
	query := fmt.Sprintf(`{
		"sort": [
			{ "date": { "order": "desc" } }
		],
		"query": {
			"match": {
				"title": "%s"
			}
		}
	}`, title)

	searchResult, err := c.Client.Search(
		c.Client.Search.WithIndex("kubesy"),
		c.Client.Search.WithBody(strings.NewReader(query)),
	) // 쿼리 실행
	if err != nil {
		return "", err
	}

	if searchResult.IsError() {
		err := fmt.Errorf("error: %s", searchResult.String())
		return "", err
	}

	// decode the response
	var response map[string]interface{}

	if err := json.NewDecoder(searchResult.Body).Decode(&response); err != nil {
		return "", err
	}

	hits := response["hits"].(map[string]interface{})["hits"].([]interface{})
	id := hits[0].(map[string]interface{})["_id"].(string)

	return id, nil

}

func (c *Client) GetPost(id string) (interface{}, error) {

	response, err := c.Client.Get("kubesy", id)
	if err != nil {
		return nil, err
	}

	// Check response status
	if response.IsError() {
		err := fmt.Errorf("error: %s", response.String())
		return nil, err
	}

	var post interface{}
	if err := json.NewDecoder(response.Body).Decode(&post); err != nil {
		return nil, err
	}

	return post, nil
}

type AdjacentPost struct {
	Title           string `json:"title"`
	PaginatorNumber int    `json:"number"`
}

func (c *Client) GetPrevNextPost(title string, row int) (AdjacentPost, AdjacentPost, int, error) {

	// 검색 쿼리 실행
	query := `{
		"sort": [
			{ "date": { "order": "desc" } }
		],
		"query": {
			"match_all": {}
		}
	}`

	searchResult, err := c.Client.Search(
		c.Client.Search.WithIndex("kubesy"),
		c.Client.Search.WithBody(strings.NewReader(query)),
	) // 쿼리 실행

	if err != nil {
		return AdjacentPost{}, AdjacentPost{}, 0, err
	}

	if searchResult.IsError() {
		err := fmt.Errorf("error: %s", searchResult.String())
		return AdjacentPost{}, AdjacentPost{}, 0, err
	}

	var response map[string]interface{}

	if err := json.NewDecoder(searchResult.Body).Decode(&response); err != nil {
		return AdjacentPost{}, AdjacentPost{}, 0, err
	}

	total := response["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)
	hits := response["hits"].(map[string]interface{})["hits"].([]interface{})

	// hits에서 현재 포스트의 인덱스 찾기

	var prev, next AdjacentPost

	var index int

	for i, hit := range hits {
		if hit.(map[string]interface{})["_source"].(map[string]interface{})["title"].(string) == title {
			index = i
			if i == 0 {
				prev.Title = ""
				prev.PaginatorNumber = 0
				next.Title = hits[i+1].(map[string]interface{})["_source"].(map[string]interface{})["title"].(string)
				next.PaginatorNumber = (i + 1) / row
			} else if i == int(total)-1 {
				prev.Title = hits[i-1].(map[string]interface{})["_source"].(map[string]interface{})["title"].(string)
				prev.PaginatorNumber = (i - 1) / row
				next.Title = ""
				next.PaginatorNumber = int(total) / row
			} else {
				prev.Title = hits[i-1].(map[string]interface{})["_source"].(map[string]interface{})["title"].(string)
				prev.PaginatorNumber = (i - 1) / row
				next.Title = hits[i+1].(map[string]interface{})["_source"].(map[string]interface{})["title"].(string)
				next.PaginatorNumber = (i + 1) / row
			}
			break
		}
	}
	return prev, next, index, nil
}

// fixme: update post
func (c *Client) UpdatePost(post Post) error {

	data, _ := json.Marshal(post)
	response, err := c.Client.Index("kubesy", bytes.NewReader(data))
	if err != nil {
		return err
	}

	// Check response status
	if response.IsError() {
		err := fmt.Errorf("error: %s", response.String())
		return err
	}

	return nil
}

func (c *Client) CountAllPosts() (int, error) {

	// 인덱스 이름 설정
	indexName := "kubesy"

	// Count API 요청
	res, err := c.Client.Count(
		c.Client.Count.WithIndex(indexName),
		c.Client.Count.WithContext(context.Background()),
	)
	if err != nil {
		log.Fatalf("Error getting document count: %s", err)
	}
	defer res.Body.Close()

	// 응답이 에러인지 확인
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return 0, err
		} else {
			return 0, err
		}
	}

	// 응답 파싱
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	// 문서 수 출력
	count := int(r["count"].(float64))

	return int(count), nil
}

type User struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	Role         string `json:"role,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

func (c *Client) GetUserByEmail(email string) (User, error) {
	// 검색 쿼리 실행
	query := fmt.Sprintf(`{
		"query": {
			"term": {
				"email": "%s"
			}
		}
	}`, email)

	searchResult, err := c.Client.Search(
		c.Client.Search.WithIndex("users"),
		c.Client.Search.WithBody(strings.NewReader(query)),
	) // 쿼리 실행

	if err != nil {
		return User{}, err
	}
	defer searchResult.Body.Close()

	if searchResult.IsError() {
		return User{}, fmt.Errorf("error: search result is error")
	}

	// 결과 파싱
	var response SearchResult
	if err := json.NewDecoder(searchResult.Body).Decode(&response); err != nil {
		return User{}, err
	}

	// 검색 결과 확인
	totalHits := response.Hits.Total.Value
	if totalHits == 0 {
		return User{}, fmt.Errorf("error: user not found")
	}
	if totalHits > 1 {
		return User{}, fmt.Errorf("error: multiple users found")
	}

	// 단일 사용자 반환
	userMap := response.Hits.Hits[0].Source

	userData, err := json.Marshal(userMap)
	if err != nil {
		return User{}, err
	}

	var user User
	if err := json.Unmarshal(userData, &user); err != nil {
		return User{}, err
	}

	return user, nil
}

func (c *Client) SaveToken(email, accessToken, refreshToken string) error {

	// index 에서 email로 유저 검색
	user, err := c.GetUserByEmail(email)
	if err != nil {
		return err
	}

	// 업데이트할 사용자 데이터 생성
	user.AccessToken = accessToken
	user.RefreshToken = refreshToken

	// 업데이트할 사용자 데이터 인덱싱
	data, _ := json.Marshal(user)
	response, err := c.Client.Index(
		"users",
		bytes.NewReader(data),
		c.Client.Index.WithDocumentID(email), // 이메일을 문서 ID로 사용
		c.Client.Index.WithRefresh("true"),   // 인덱스 즉시 갱신
	)
	if err != nil {
		return err
	}

	// Check response status
	if response.IsError() {
		err := fmt.Errorf("error: %s", response.String())
		return err
	}

	return nil
}

func (c *Client) DeleteToken(email string) error {
	// index에서 email로 사용자 검색
	user, err := c.GetUserByEmail(email)
	if err != nil {
		return err
	}

	// 업데이트할 사용자 데이터 생성
	user.AccessToken = ""
	user.RefreshToken = ""

	// 업데이트할 사용자 데이터 인덱싱
	data, _ := json.Marshal(user)
	response, err := c.Client.Index(
		"users",
		bytes.NewReader(data),
		c.Client.Index.WithDocumentID(email), // 이메일을 문서 ID로 사용
		c.Client.Index.WithRefresh("true"),   // 인덱스 즉시 갱신
	)
	if err != nil {
		return err
	}

	// Check response status
	if response.IsError() {
		err := fmt.Errorf("error: %s", response.String())
		return err
	}

	return nil
}
