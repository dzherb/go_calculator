package repo

import (
	"context"
	"time"

	"github.com/dzherb/go_calculator/calculator/internal/storage"
	"github.com/georgysavva/scany/v2/pgxscan"
)

type User struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"          db:"password_hash"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	Get(id uint64) (User, error)
	GetByCredentials(username, password string) (User, error)
	Create(user User) (User, error)
}

type UserRepositoryImpl struct {
	db storage.Connection
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{
		db: storage.Conn(),
	}
}

func (ur *UserRepositoryImpl) Get(id uint64) (User, error) {
	user := User{}
	err := pgxscan.Get(
		context.Background(),
		ur.db,
		&user,
		`SELECT id, username, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1;`,
		id,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (ur *UserRepositoryImpl) Create(user User) (User, error) {
	err := pgxscan.Get(
		context.Background(),
		ur.db,
		&user,
		`INSERT INTO users (username, password_hash)
		VALUES ($1, crypt($2, gen_salt('bf')))
		RETURNING id, username, password_hash, created_at, updated_at;`,
		user.Username, user.Password,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (ur *UserRepositoryImpl) GetByCredentials(
	username, password string,
) (User, error) {
	user := User{}
	err := pgxscan.Get(
		context.Background(),
		ur.db,
		&user,
		`SELECT id, username, password_hash, created_at, updated_at
		FROM users
		WHERE username = $1 AND crypt($2, password_hash) = password_hash;`,
		username, password,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}
