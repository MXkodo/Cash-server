package repo

import (
	"context"
	"fmt"

	"github.com/MXkodo/cash-server/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DocRepo struct {
	db *pgxpool.Pool
}

func NewDocRepo(db *pgxpool.Pool) *DocRepo {
	return &DocRepo{db}
}

func (r *DocRepo) CreateDocument(ctx context.Context, userId int, name, mime, filePath string, isPublic bool) (int, error) {
	var id int
	err := r.db.QueryRow(
		ctx,
		`INSERT INTO documents (user_id, name, mime, file_path, is_public)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		userId, name, mime, filePath, isPublic,
	).Scan(&id)
	return id, err
}

func (r *DocRepo) GetDocuments(ctx context.Context, userId int, login, key, value string, limit string) ([]models.Document, error) {
	var docs []models.Document
	query := `SELECT id, user_id, name, mime, file_path, is_public, created_at 
			  FROM documents WHERE user_id = $1`

	var args []interface{}
	args = append(args, userId)

	if login != "" {
		query += " AND user_login = $2" 
		args = append(args, login)
	}

	if key != "" && value != "" {
		query += fmt.Sprintf(" AND %s = $%d", key, len(args)+1)
		args = append(args, value)
	}

	query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
	args = append(args, limit)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var doc models.Document
		if err := rows.Scan(&doc.ID, &doc.UserID, &doc.Name, &doc.Mime, &doc.FilePath, &doc.IsPublic, &doc.CreatedAt); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func (r *DocRepo) GetDocumentByID(ctx context.Context, id int) (models.Document, error) {
	var doc models.Document
	err := r.db.QueryRow(
		ctx,
		`SELECT id, user_id, name, mime, file_path, is_public, created_at 
		 FROM documents WHERE id = $1`, id,
	).Scan(&doc.ID, &doc.UserID, &doc.Name, &doc.Mime, &doc.FilePath, &doc.IsPublic, &doc.CreatedAt)
	return doc, err
}

func (r *DocRepo) DeleteDocument(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM documents WHERE id = $1`, id)
	return err
}
