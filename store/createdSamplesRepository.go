package store

import (
	"DriveApi/internal/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type CreatedSamplesRepository struct {
	store *Store
}

type Creator interface {
	/*Create(c *model.CreatedSample) error
	CreateTx(tx *sqlx.Tx, c *model.CreatedSample) error
	GetUserAmount(userID int) (int, error)
	GetUserAmountTx(tx *sqlx.Tx, userID int) (int, error)
	GetUserAll(userID int) ([]model.Sample, error)
	GetUserAllTx(tx *sqlx.Tx, userID int) ([]model.Sample, error)*/
}

func (cr *CreatedSamplesRepository) Create(c *model.CreatedSample) error {
	return cr.CreateTx(nil, c)
}

func (cr *CreatedSamplesRepository) CreateTx(tx *sqlx.Tx, c *model.CreatedSample) error {
	if err := cr.store.db.QueryRow(
		tx,
		`INSERT INTO created_samples (sample_id, user_id) VALUES ($1, $2);`,
		c.SampleID, c.UserID,
	).Err(); err != nil {
		return fmt.Errorf("failed to insert data to table samples; %w", err)
	}

	return nil
}

func (cr *CreatedSamplesRepository) GetUserAmount(userID int) (int, error) {
	return cr.GetUserAmountTx(nil, userID)
}

func (cr *CreatedSamplesRepository) GetUserAmountTx(tx *sqlx.Tx, userID int) (int, error) {
	var res int
	row := cr.store.db.QueryRow(tx, `SELECT COUNT(*) FROM created_samples WHERE user_id=$1;`, userID)

	err := row.Scan(&res)
	if err != nil {
		return -1, fmt.Errorf("failed to scan amount of created samples %w", err)
	}

	return res, nil
}

func (cr *CreatedSamplesRepository) GetUserAll(userID int) ([]model.Sample, error) {
	return cr.GetUserAllTx(nil, userID)
}

func (cr *CreatedSamplesRepository) GetUserAllTx(tx *sqlx.Tx, userID int) ([]model.Sample, error) {
	samples := make([]model.Sample, 0)
	samplesID := make([]int, 0)
	rows, err := cr.store.db.Query(tx, `SELECT sample_id FROM created_samples WHERE user_id=$1;`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to select from table samples; %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var sampleID int
		err := rows.Scan(&sampleID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sampleID; %w", err)
		}
		samplesID = append(samplesID, sampleID)
	}

	return samples, nil
}
