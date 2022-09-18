package store

import (
	"DriveApi/internal/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type LikedSamplesRepository struct {
	store *Store
}

type Liker interface {
	Create(c *model.LikedSamples) error
	CreateTx(tx *sqlx.Tx, c *model.LikedSamples) error
	GetUserAmount(userID int) (int, error)
	GetUserAmountTx(tx *sqlx.Tx, userID int) (int, error)
	GetUserAll(userID int) ([]model.Sample, error)
	GetUserAllTx(tx *sqlx.Tx, userID int) ([]model.Sample, error)
}

func (l *LikedSamplesRepository) Create(c *model.LikedSamples) error {
	return l.CreateTx(nil, c)
}

func (l *LikedSamplesRepository) CreateTx(tx *sqlx.Tx, c *model.LikedSamples) error {
	if err := l.store.db.QueryRow(
		tx,
		`INSERT INTO liked_samples (sample_id, user_id) VALUES ($1, $2);`,
		c.SampleID, c.UserID,
	).Err(); err != nil {
		return fmt.Errorf("failed to insert data to table samples; %w", err)
	}

	return nil
}

func (l *LikedSamplesRepository) GetUserAmount(userID int) (int, error) {
	return l.GetUserAmountTx(nil, userID)
}

func (l *LikedSamplesRepository) GetUserAmountTx(tx *sqlx.Tx, userID int) (int, error) {
	var res int
	row := l.store.db.QueryRow(tx, `SELECT COUNT(*) FROM liked_samples WHERE user_id=$1;`, userID)

	err := row.Scan(&res)
	if err != nil {
		return -1, fmt.Errorf("failed to scan amount of created samples %w", err)
	}

	return res, nil
}

func (l *LikedSamplesRepository) GetUserAll(userID int) ([]model.Sample, error) {
	return l.GetUserAllTx(nil, userID)
}

func (l *LikedSamplesRepository) GetUserAllTx(tx *sqlx.Tx, userID int) ([]model.Sample, error) {
	samplesID := make([]int, 0)
	rows, err := l.store.db.Query(tx, `SELECT DISTINCT sample_id, user_id FROM liked_samples WHERE user_id=$1;`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to select from table samples; %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		sampleID := &model.LikedSamples{}
		err := rows.StructScan(&sampleID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sampleID; %w", err)
		}
		samplesID = append(samplesID, sampleID.SampleID)
	}

	samples := make([]model.Sample, 0)
	for _, v := range samplesID {
		s := model.Sample{}
		row := l.store.db.QueryRow(tx, `SELECT id, name, path, author, author_id, type FROM samples WHERE id=$1`, v)
		err := row.StructScan(&s)
		if err != nil {
			return nil, fmt.Errorf("failed to struct scan sample %w", err)
		}
		samples = append(samples, s)
	}

	return samples, nil
}
