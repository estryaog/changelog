package database

import (
	"database/sql"

	"github.com/estryaog/changelog/internal/types"
	"github.com/gin-gonic/gin"
)

type ChangelogStore interface {
	CreateChangelog(ctx *gin.Context, changelog *types.Changelog) (*types.Changelog, error)
	DeleteChangelog(ctx *gin.Context, id string) error
	UpdateChangelog(ctx *gin.Context, id string, changelog *types.Changelog) (*types.Changelog, error)
	GetChangelog(ctx *gin.Context, id string) (*types.Changelog, error)
	GetChangelogs(ctx *gin.Context, limit int, orderBy string) ([]*types.Changelog, error)
}

type PostgresChangelogsStore struct {
	db *sql.DB
}

func NewPostgresChangelogStore(db *ServiceImpl) ChangelogStore {
	return &PostgresChangelogsStore{db: db.db}
}

func (db PostgresChangelogsStore) CreateChangelog(ctx *gin.Context, changelog *types.Changelog) (*types.Changelog, error) {
	query := `INSERT INTO changelogs (version, content) VALUES ($1, $2) RETURNING id, created_at`
	err := db.db.QueryRowContext(ctx, query, changelog.Version, changelog.Content).Scan(&changelog.Id, &changelog.CreatedAt)
	if err != nil {
		return nil, err
	}

	return changelog, nil
}

func (db PostgresChangelogsStore) DeleteChangelog(ctx *gin.Context, id string) error {
	query := `DELETE FROM changelogs WHERE id = $1`
	_, err := db.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (db PostgresChangelogsStore) UpdateChangelog(ctx *gin.Context, id string, changelog *types.Changelog) (*types.Changelog, error) {
	query := `UPDATE changelogs SET version = $1, content = $2 WHERE id = $3 RETURNING created_at`
	err := db.db.QueryRowContext(ctx, query, changelog.Version, changelog.Content, id).Scan(&changelog.CreatedAt)
	if err != nil {
		return nil, err
	}

	changelog.Id = id
	return changelog, nil
}

func (db PostgresChangelogsStore) GetChangelog(ctx *gin.Context, id string) (*types.Changelog, error) {
	query := `SELECT id, version, content, created_at FROM changelogs WHERE id = $1`
	changelog := &types.Changelog{}
	err := db.db.QueryRowContext(ctx, query, id).Scan(&changelog.Id, &changelog.Version, &changelog.Content, &changelog.CreatedAt)
	if err != nil {
		return nil, err
	}

	return changelog, nil
}

func (db PostgresChangelogsStore) GetChangelogs(ctx *gin.Context, limit int, orderBy string) ([]*types.Changelog, error) {
	query := `SELECT id, version, content, created_at FROM changelogs ORDER BY $1 DESC LIMIT $2`
	rows, err := db.db.QueryContext(ctx, query, orderBy, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	changelogs := []*types.Changelog{}
	for rows.Next() {
		changelog := &types.Changelog{}
		err := rows.Scan(&changelog.Id, &changelog.Version, &changelog.Content, &changelog.CreatedAt)
		if err != nil {
			return nil, err
		}

		changelogs = append(changelogs, changelog)
	}

	return changelogs, nil
}
