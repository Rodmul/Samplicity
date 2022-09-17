package store

import (
	"DriveApi/internal/model"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	salt       = "hjqrhjqw124617ajfhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type UserRepository struct {
	store *Store
}

type tokenClaims struct {
	jwt.StandardClaims
	UserID int `json:"userID"`
}

type Profiler interface {
	Create(c *model.User) error
	CreateTx(tx *sqlx.Tx, c *model.User) error
	GetUser(username, password string) (*model.User, error)
	GetUserTx(tx *sqlx.Tx, username, password string) (*model.User, error)
	GetUserByID(id int) (*model.User, error)
	GetUserByIDTx(tx *sqlx.Tx, id int) (*model.User, error)
	GenerateToken(username, password string) (string, error)
	GetPasswordHash(passport string) string
	ParseToken(token string) (int, error)
}

func (r *UserRepository) Create(c *model.User) error {
	return r.CreateTx(nil, c)
}

func (r *UserRepository) CreateTx(tx *sqlx.Tx, c *model.User) error {

	hashPassport := r.GetPasswordHash(c.Password)
	if err := r.store.db.QueryRow(
		tx,
		`INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3);`,
		c.Username, hashPassport, c.Email,
	).Err(); err != nil {
		return fmt.Errorf("failed to insert data to table samples; %w", err)
	}

	return nil
}

func (r *UserRepository) GetUser(username, password string) (*model.User, error) {
	return r.GetUserTx(nil, username, password)
}

func (r *UserRepository) GetUserTx(tx *sqlx.Tx, username, password string) (*model.User, error) {
	user := &model.User{}
	row := r.store.db.QueryRow(tx, `SELECT id, username, email, password_hash FROM users WHERE username=$1 AND password_hash=$2`, username, password)

	err := row.StructScan(user)
	if err != nil {
		return nil, fmt.Errorf("failed to struct scan user; %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByID(id int) (*model.User, error) {
	return r.GetUserByIDTx(nil, id)
}

func (r *UserRepository) GetUserByIDTx(tx *sqlx.Tx, id int) (*model.User, error) {
	user := &model.User{}
	row := r.store.db.QueryRow(tx, `SELECT id, username, email FROM users WHERE id=$1`, id)

	err := row.StructScan(user)
	if err != nil {
		return nil, fmt.Errorf("failed to struct scan user; %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetPasswordHash(passport string) string {
	hash := sha1.New()
	hash.Write([]byte(passport))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (r *UserRepository) GenerateToken(username, password string) (string, error) {
	user, err := r.GetUser(username, r.GetPasswordHash(password))
	if err != nil {
		return "", fmt.Errorf("failed to get user; %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})

	return token.SignedString([]byte(signingKey))
}

func (r *UserRepository) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserID, nil
}
