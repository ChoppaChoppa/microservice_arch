package PgDataBase

import (
	"context"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"sub_cache/Models"
)

type DB struct {
	conn *sqlx.DB
}

func Connection(uri string) (*DB, error) {
	connUri := uri //"postgresql://selectel:selectel@192.168.3.30:5432/"

	dataBase, err := sqlx.Connect("pgx", connUri)
	if err != nil {
		return nil, fmt.Errorf("failed connect: %w", err)
	}

	dataBase.SetMaxOpenConns(1)
	dataBase.SetMaxIdleConns(3)

	//driver, errDriver := postgres.WithInstance(dataBase.DB, &postgres.Config{
	// DatabaseName: "postgres",
	// SchemaName: "public",
	//})
	//if errDriver != nil {
	// return nil, fmt.Errorf("migrate instance: %w", errDriver)
	//}
	//
	//m, err := migrate.NewWithDatabaseInstance(
	// "file://migrations",
	// "postgres", driver)
	//if err != nil {
	// return nil, fmt.Errorf("migrate: %w", err)
	//}
	//if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
	// return nil, fmt.Errorf("migrate up: %w", err)
	//}
	//
	return &DB{
		conn: dataBase,
	}, nil
}
//TODO проверить можно возвращать боле значений в постгрес

func (dataBase *DB) Add(ctx context.Context, user Models.) (Models.OrderInfo, error) {
	var id string

	hash, errHash := HashPassword(user.Password)
	if errHash != nil {
		fmt.Errorf("failed to generate hash %v", errHash)
		return Models.User{}, errHash
	}

	if errCreate := dataBase.conn.GetContext(ctx, &id,
		`INSERT INTO users(login, password)
				VALUES ($1, $2)
				RETURNING id`, user.Login, hash); errCreate != nil {
		fmt.Errorf("failed to create user %v", errCreate.Error())
		return Models.User{}, errCreate
	}

	User := Models.User{
		ID:       id,
		Login:    user.Login,
		Password: hash,
	}

	return User, nil
}

func (dataBase *DB) Get(ctx context.Context, id string) (Models.User, error) {
	var user Models.User

	if errGet := dataBase.conn.GetContext(ctx, &user,
		`SELECT id, login
FROM users`); errGet != nil {
		fmt.Errorf("get by id: %v", errGet)
		return Models.User{}, errGet
	}

	return user, nil
}

func HashPassword(password string) (string, error) {
	hash, errGenerateHash := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errGenerateHash != nil {
		return "", errGenerateHash
	}

	return string(hash), nil
}
