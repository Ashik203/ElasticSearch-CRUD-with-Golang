package post

import (
	"context"
	"elasticsearch/internal/pkg/storage"
	"time"

	"github.com/google/uuid"
)

type service struct {
	storage storage.PostStorer
}

func (s service) create(ctx context.Context, req createRequest) (createResponse, error) {
	id := uuid.New().String()
	cr := time.Now().UTC()
	doc := storage.Post{
		ID:        id,
		Title:     req.Title,
		Text:      req.Text,
		Tags:      req.Tags,
		CreatedAt: &cr,
	}

	if err := s.storage.Insert(ctx, doc); err != nil {
		return createResponse{}, err
	}

	return createResponse{ID: id, Title: req.Title, Text: req.Text,
		Tags: req.Tags}, nil
}

func (s service) update(ctx context.Context, req updateRequest) (updateResponse, error) {
	doc := storage.Post{
		ID:    req.ID,
		Title: req.Title,
		Text:  req.Text,
		Tags:  req.Tags,
	}

	if err := s.storage.Update(ctx, doc); err != nil {
		return updateResponse{}, err
	}
	return updateResponse{Title: req.Title, Text: req.Text,
		Tags: req.Tags}, nil
}

func (s service) delete(ctx context.Context, req deleteRequest) (deleteResponse, error) {
	if err := s.storage.Delete(ctx, req.ID); err != nil {
		return deleteResponse{}, err
	}
	return deleteResponse{ID: req.ID}, nil
}

func (s service) find(ctx context.Context, req findRequest) (findResponse, error) {
	post, err := s.storage.FindOne(ctx, req.ID)
	if err != nil {
		return findResponse{}, err
	}

	return findResponse{
		ID:        post.ID,
		Title:     post.Title,
		Text:      post.Text,
		Tags:      post.Tags,
		CreatedAt: *post.CreatedAt,
	}, nil
}
