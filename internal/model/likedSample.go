package model

type LikedSamples struct {
	SampleID int `db:"sample_id"`
	UserID   int `db:"user_id"`
}
