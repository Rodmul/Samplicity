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
	Create(c *model.Sample) (int, error)
	CreateTx(tx *sqlx.Tx, c *model.Sample) (int, error)
	GetAll() ([]model.Sample, error)
	GetAllTx(tx *sqlx.Tx) ([]model.Sample, error)
	GetByName(name string) (*model.Sample, error)
	GetByNameTx(tx *sqlx.Tx, name string) (*model.Sample, error)
	GetUserCreated(userID int) ([]model.Sample, error)
	GetUserCreatedTx(tx *sqlx.Tx, userID int) ([]model.Sample, error)
	GetUserCreatedAmount(userID int) (int, error)
	GetUserCreatedAmountTx(tx *sqlx.Tx, userID int) (int, error)
}

func (r *SampleRepository) Create(c *model.Sample) (int, error) {
	return r.CreateTx(nil, c)
}

func (r *SampleRepository) CreateTx(tx *sqlx.Tx, c *model.Sample) (int, error) {
	var i int

	if err := r.store.db.QueryRow(
		tx,
		`INSERT INTO samples (name, path, author, author_id, type) VALUES ($1, $2, $3, $4, $5);`,
		c.Name, c.Path, c.Author, c.AuthorID, c.Type,
	).Err(); err != nil {
		return -1, fmt.Errorf("failed to insert data to table samples; %w", err)
	}

	row := r.store.db.QueryRow(tx, `SELECT id FROM samples WHERE name=$1`, c.Name)
	err := row.Scan(&i)
	if err != nil {
		return -1, fmt.Errorf("failed to struct scan sample %w", err)
	}

	return i, nil
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
	row := r.store.db.QueryRow(tx, `SELECT id, name, path, author, author_id, type FROM samples WHERE name=$1`, name)

	err := row.StructScan(&sample)
	if err != nil {
		return nil, fmt.Errorf("failed to struct scan sample %w", err)
	}

	return &sample, nil
}

func (r *SampleRepository) GetUserCreated(userID int) ([]model.Sample, error) {
	return r.GetUserCreatedTx(nil, userID)
}

func (r *SampleRepository) GetUserCreatedTx(tx *sqlx.Tx, userID int) ([]model.Sample, error) {
	samples := make([]model.Sample, 0)
	rows, err := r.store.db.Query(tx, `SELECT id, name, path, author, author_id, type FROM samples WHERE author_id=$1`, userID)
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

func (r *SampleRepository) GetUserCreatedAmount(userID int) (int, error) {
	return r.GetUserCreatedAmountTx(nil, userID)
}

func (r *SampleRepository) GetUserCreatedAmountTx(tx *sqlx.Tx, userID int) (int, error) {
	var res int
	row := r.store.db.QueryRow(tx, `SELECT COUNT(*) FROM samples WHERE author_id=$1;`, userID)

	err := row.Scan(&res)
	if err != nil {
		return -1, fmt.Errorf("failed to scan amount of created samples %w", err)
	}

	return res, nil
}
