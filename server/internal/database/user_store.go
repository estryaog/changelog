package database

import (
	"database/sql"

	"github.com/estryaog/changelog/internal/types"
	"github.com/gin-gonic/gin"
)

type UserStore interface {
	CreateUser(ctx *gin.Context, user *types.User) (*types.User, error)
	DeleteUser(ctx *gin.Context, id string) error
	GetUser(ctx *gin.Context, key, value string) (*types.User, error)
}

type PostgresUsersStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *ServiceImpl) UserStore {
	return &PostgresUsersStore{db: db.db}
}

func (db PostgresUsersStore) CreateUser(ctx *gin.Context, user *types.User) (*types.User, error) {
	query := `INSERT INTO users (id, email, password, is_admin) VALUES ($1, $2, $3, $4) RETURNING id`
	err := db.db.QueryRowContext(ctx, query, user.Id, user.Email, user.Password, user.IsAdmin).Scan(&user.Id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db PostgresUsersStore) DeleteUser(ctx *gin.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (db PostgresUsersStore) GetUser(ctx *gin.Context, key, value string) (*types.User, error) {
	query := `SELECT id, email, password, is_admin FROM users WHERE ` + key + ` = $1`
	user := &types.User{}
	err := db.db.QueryRowContext(ctx, query, value).Scan(&user.Id, &user.Email, &user.Password, &user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return user, nil
}
