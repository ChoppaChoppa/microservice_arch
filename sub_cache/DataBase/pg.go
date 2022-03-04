package DataBase

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"sub_cache/Models"
)

type DB struct {
	connect *sqlx.DB
}

func Connection(url string) (*DB, error) {
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

func (dataBase *DB) GetOrders(ctx context.Context, id string) (Models.OrderInfo, error) {
	var order Models.OrderInfo

	if errGetOrder := dataBase.connect.GetContext(ctx, &order, `
			SELECT order_uid, track_number, entry, locale, internal_signature,
                        customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard
			FROM order_info
			WHERE id = $1
		`, id); errGetOrder != nil {
		return Models.OrderInfo{}, fmt.Errorf("get order: %v", errGetOrder)
	}
	order.ID = id

	if errGetDelivery := dataBase.connect.GetContext(ctx, &order.Delivery, `
	SELECT name, phone, zip, city, address, region, email
	FROM deliver
	WHERE order_id = $1
	`, order.ID); errGetDelivery != nil {
		return Models.OrderInfo{}, fmt.Errorf("get delivery: %v", errGetDelivery)
	}

	if errGetPayments := dataBase.connect.GetContext(ctx, &order.Payment, `
	SELECT transaction, request_id, currency, provider, amount,
                    payment_dt, bank, delivery_cost, goods_total, custom_fee
	FROM payment
	WHERE order_id = $1
	`, order.ID); errGetPayments != nil {
		return Models.OrderInfo{}, fmt.Errorf("get payments: %v", errGetPayments)
	}

	if errGetItems := dataBase.connect.SelectContext(ctx, &order.Items, `
	SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
	FROM item
	WHERE order_id = $1
	`, order.ID); errGetItems != nil {
		return Models.OrderInfo{}, fmt.Errorf("get items: %v", errGetItems)
	}

	return order, nil
}

func (dataBase *DB) GetLasts(ctx context.Context, count int) ([]Models.OrderInfo, error){
	var Orders []Models.OrderInfo

	fmt.Println("sstart get lasts")
	if errGetOrder := dataBase.connect.SelectContext(ctx, &Orders, `
			SELECT id, order_uid, track_number, entry, locale, internal_signature,
                        customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard
			FROM order_info
			ORDER BY id DESC LIMIT $1`, count);
	errGetOrder != nil {
		fmt.Println("get order ", errGetOrder)
		return nil, fmt.Errorf("get order: %v", errGetOrder)
	}

	fmt.Println("orders got", Orders)

	for i, _ := range Orders {
		if errGetDelivery := dataBase.connect.GetContext(ctx, &Orders[i].Delivery, `
			SELECT name, phone, zip, city, address, region, email
			FROM deliver
			WHERE order_id = $1`, Orders[i].ID);
		errGetDelivery != nil {
			fmt.Println("get delivery ", errGetDelivery)
			return nil, fmt.Errorf("get delivery: %v", errGetDelivery)
		}

		fmt.Println("delivery got")

		if errGetPayments := dataBase.connect.GetContext(ctx, &Orders[i].Payment, `
			SELECT transaction, request_id, currency, provider, amount,
                    payment_dt, bank, delivery_cost, goods_total, custom_fee
			FROM payment
			WHERE order_id = $1`, Orders[i].ID);
		errGetPayments != nil {
			fmt.Println("get payments ", errGetPayments)
			return nil, fmt.Errorf("get payments: %v", errGetPayments)
		}

		fmt.Println("payments got ", Orders[i].Payment)

		if errGetItems := dataBase.connect.SelectContext(ctx, &Orders[i].Items, `
			SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
			FROM item
			WHERE order_id = $1`, Orders[i].ID);
		errGetItems != nil {
			fmt.Println("get items ", errGetItems)
			return nil, fmt.Errorf("get items: %v", errGetItems)
		}

		fmt.Println("items got ", Orders[i].Items)
	}

	return Orders, nil
}