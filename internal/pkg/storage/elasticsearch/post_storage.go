package elasticsearch

import (
	"bytes"
	"context"
	"elasticsearch/internal/pkg/domain"
	"elasticsearch/internal/pkg/storage"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var _ storage.PostStorer = PostStorage{}

type PostStorage struct {
	elastic ElasticSearch
	timeout time.Duration
}

func (p PostStorage) Update(ctx context.Context, post storage.Post) error {
	panic("unimplemented")
}

func NewPostStorage(elastic ElasticSearch) (PostStorage, error) {
	return PostStorage{
		elastic: elastic,
		timeout: time.Second * 10,
	}, nil
}

func (p PostStorage) Insert(ctx context.Context, post storage.Post) error {
	bdy, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("insert: mershal:%w", err)
	}

	req := esapi.CreateRequest{
		Index:      p.elastic.alias,
		DocumentID: post.ID,
		Body:       bytes.NewReader(bdy),
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := req.Do(ctx, p.elastic.client)
	if err != nil {
		return fmt.Errorf("insert: request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode == 409 {
		return domain.ErrConflict
	}

	if res.IsError() {
		return fmt.Errorf("insert: response: %s", res.String())
	}

	return nil
}

func (p PostStorage) Delete(ctx context.Context, id string) error {
	req := esapi.DeleteRequest{
		Index:      p.elastic.alias,
		DocumentID: id,
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := req.Do(ctx, p.elastic.client)
	if err != nil {
		return fmt.Errorf("delete: request:%w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return domain.ErrNotFound
	}

	if res.IsError() {
		return fmt.Errorf("delete: response:%s", res.String())
	}
	return nil
}

func (p PostStorage) FindOne(ctx context.Context, id string) (storage.Post, error) {
	req := esapi.GetRequest{
		Index:      p.elastic.alias,
		DocumentID: id,
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := req.Do(ctx, p.elastic.client)
	if err != nil {
		return storage.Post{}, fmt.Errorf("find one : request:%w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return storage.Post{}, domain.ErrNotFound
	}

	if res.IsError() {
		return storage.Post{}, fmt.Errorf("find one: response:%s", res.String())
	}

	var (
		post storage.Post
		body document
	)

	body.Source = &post
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return storage.Post{}, fmt.Errorf("find one: decode:%w", err)
	}

	return post, nil
}
