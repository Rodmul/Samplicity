package store

import (
	"DriveApi/internal/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type SampleRepository struct {
	store *Store
}

type Sampler interface {
	Create(c *model.Sample) error
	CreateTx(tx *sqlx.Tx, c *model.Sample) error
	GetAll() ([]model.Sample, error)
	GetAllTx(tx *sqlx.Tx) ([]model.Sample, error)
	GetByName(name string) (*model.Sample, error)
	GetByNameTx(tx *sqlx.Tx, name string) (*model.Sample, error)
}

func (r *SampleRepository) Create(c *model.Sample) error {
	return r.CreateTx(nil, c)
}

func (r *SampleRepository) CreateTx(tx *sqlx.Tx, c *model.Sample) error {
	if err := r.store.db.QueryRow(
		tx,
		`INSERT INTO samples (name, path, author, type) VALUES ($1, $2, $3, $4);`,
		c.Name, c.Path, c.Author, c.Type,
	).Err(); err != nil {
		return fmt.Errorf("failed to insert data to table samples; %w", err)
	}

	return nil
}

func (r *SampleRepository) GetAll() ([]model.Sample, error) {
	return r.GetAllTx(nil)
}

func (r *SampleRepository) GetAllTx(tx *sqlx.Tx) ([]model.Sample, error) {
	samples := make([]model.Sample, 0)
	rows, err := r.store.db.Query(tx, `SELECT id, name, path, author, type FROM samples;`)
	if err != nil {
		return nil, fmt.Errorf("failed to select from table samples; %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		sample := model.Sample{}
		err := rows.StructScan(&sample)
		if err != nil {
			return nil, fmt.Errorf("failed to scan structure; %w", err)
		}
		samples = append(samples, sample)
	}

	return samples, nil
}

func (r *SampleRepository) GetByName(name string) (*model.Sample, error) {
	return r.GetByNameTx(nil, name)
}

func (r *SampleRepository) GetByNameTx(tx *sqlx.Tx, name string) (*model.Sample, error) {
	sample := model.Sample{}
	row := r.store.db.QueryRow(tx, `SELECT id, name, path, author, type FROM samples WHERE name=$1`, name)

	err := row.StructScan(&sample)
	if err != nil {
		return nil, fmt.Errorf("failed to struct scan sample %w", err)
	}

	return &sample, nil
}
