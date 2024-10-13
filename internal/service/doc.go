package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/MXkodo/cash-server/internal/repo"
	"github.com/MXkodo/cash-server/models"
	"github.com/redis/go-redis/v9"
)

type DocService struct {
	docRepo *repo.DocRepo
	rdb     *redis.Client
}

func NewDocService(docRepo *repo.DocRepo, rdb *redis.Client) *DocService {
	return &DocService{docRepo, rdb}
}

func (s *DocService) UploadDocument(ctx context.Context, userId int, name string, mime string, isPublic bool, file multipart.File) (map[string]string, error) {
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, err
	}

	filePath := filepath.Join(uploadDir, name)

	out, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return nil, err
	}

	docID, err := s.docRepo.CreateDocument(ctx, userId, name, mime, filePath, isPublic)
	if err != nil {
		return nil, err
	}

	return map[string]string{"id": fmt.Sprintf("%d", docID), "file_path": filePath}, nil
}

func (s *DocService) GetDocuments(ctx context.Context, userId int, login, key, value, limit string) ([]models.Document, error) {
	cacheKey := fmt.Sprintf("documents:user:%d:login:%s:key:%s:value:%s:limit:%s", userId, login, key, value, limit)
	cachedDocs, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var docs []models.Document
		if err := json.Unmarshal([]byte(cachedDocs), &docs); err == nil {
			return docs, nil
		}
	}

	docs, err := s.docRepo.GetDocuments(ctx, userId, login, key, value, limit)
	if err != nil {
		return docs, err
	}

	cachedDocsBytes, _ := json.Marshal(docs)
	s.rdb.Set(ctx, cacheKey, cachedDocsBytes, 0)

	return docs, nil
}

func (s *DocService) GetDocument(ctx context.Context, id int) (models.Document, error) {
	cacheKey := fmt.Sprintf("document:%d", id)
	cachedDoc, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var doc models.Document
		if err := json.Unmarshal([]byte(cachedDoc), &doc); err == nil {
			return doc, nil
		}
	}

	doc, err := s.docRepo.GetDocumentByID(ctx, id)
	if err != nil {
		return doc, err
	}

	cachedDocBytes, _ := json.Marshal(doc)
	s.rdb.Set(ctx, cacheKey, cachedDocBytes, 0)

	return doc, nil
}

func (s *DocService) DeleteDocument(ctx context.Context, id int) error {
	cacheKey := fmt.Sprintf("document:%d", id)
	s.rdb.Del(ctx, cacheKey)

	return s.docRepo.DeleteDocument(ctx, id)
}
