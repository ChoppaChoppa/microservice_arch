package DataBase

import (
	"github.com/jmoiron/sqlx"
)

type DB struct {
	connect *sqlx.DB
}

func Connection(url string) (*DB, error){
	connUri := url

	dataBase, err := sqlx.Connect("pgx", connUri)
	if err != nil {
		return nil, err
	}

	dataBase.SetMaxOpenConns(1)
	dataBase.SetMaxIdleConns(3)

	//driver, errDriver := postgres.WithInstance(dataBase.DB, &postgres.Config{
	//	DatabaseName: "postgres",
	//	SchemaName: "public",
	//})
	//if errDriver != nil {
	//	return nil, fmt.Errorf("migrate instance: %w", errDriver)
	//}
	//
	//m, err := migrate.NewWithDatabaseInstance(
	//	"file://migrations",
	//	"postgres", driver)
	//if err != nil {
	//	return nil, fmt.Errorf("migrate: %w", err)
	//}
	//if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
	//	return nil, fmt.Errorf("migrate up: %w", err)
	//}

	return &DB{
		connect: dataBase,
	}, nil
}

func (dataBase *DB) GetOrders(id string){
	
}