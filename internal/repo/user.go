package repo

import (
	"context"

	"github.com/MXkodo/cash-server/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) CreateUser(login, hashedPassword string) error {
	_, err := r.db.Exec(context.Background(), "INSERT INTO users (login, password) VALUES ($1, $2)", login, hashedPassword)
	return err
}

func (r *UserRepo) FindUserByLogin(login string) (models.User, error) {
	var user models.User
	err := r.db.QueryRow(context.Background(), "SELECT id, login, password, created_at FROM users WHERE login = $1", login).
		Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
