package PgDataBase

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"sub_db/Models"
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

func (dataBase *DB) Add(ctx context.Context, order Models.OrderInfo) (Models.OrderInfo, error) {
	var deliverID, paymentsID string

	if errTx := dataBase.InTx(ctx, sql.LevelReadCommitted, func(tx sqlx.Tx) error {
		fmt.Println("start tx")

		var orderID string
		errExecOrder := tx.Get(&orderID, `
		INSERT INTO order_info (order_uid, track_number, entry, locale, internal_signature,
                     customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
		`, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
			order.CustomerID, order.DeliveryService, order.ShardKey,
			order.SmID, order.DateCreated, order.OofShard)
		if errExecOrder != nil {
			fmt.Println(errExecOrder)
			return fmt.Errorf("failed order insert: %v", errExecOrder)
		}
		order.ID = orderID

		if errExecDelivery := tx.Get(&deliverID, `
		INSERT INTO deliver(order_id, name, phone, zip, city, address, region, email)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;
		`,order.ID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City,
			order.Delivery.Address, order.Delivery.Region, order.Delivery.Email); errExecDelivery != nil {
			fmt.Println(errExecDelivery)
			return fmt.Errorf("failed delivery insert: %v", errExecDelivery)
		}
		fmt.Println("delivery end")

		if errExecPayments := tx.Get(&paymentsID, `
		INSERT INTO payment(order_id, transaction, request_id, currency, provider, amount,
                 payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id;
		`, order.ID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
			order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt,
			order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee); errExecPayments != nil {
			fmt.Println(errExecPayments)
			return fmt.Errorf("failed payment insert: %v", errExecPayments)
		}
		fmt.Println("payment end")

		for _, v := range order.Items {
			_, errExecItem := tx.Exec(`
				INSERT INTO item(order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
				`, order.ID, v.ChrtID, v.TrackNumber, v.Price, v.RID, v.Name, v.Sale, v.Size, v.TotalPrice,
				v.NmID, v.Brand, v.Status)
			if errExecItem != nil {

				return fmt.Errorf("failed insert item: %v", errExecItem)
			}
		}
		fmt.Println("item end")

		return nil
	}); errTx != nil {
		return Models.OrderInfo{}, errTx
	}

	fmt.Println("tx complete")
	return order, nil
}

func (dataBase *DB) GetOrder(ctx context.Context, id string) (Models.OrderInfo, error) {
	var order Models.OrderInfo

	if errGetOrder := dataBase.conn.GetContext(ctx, &order, `
			SELECT order_uid, track_number, entry, locale, internal_signature,
                        customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard
			FROM order_info
			WHERE id = $1
		`, id); errGetOrder != nil {
		return Models.OrderInfo{}, fmt.Errorf("get order: %v", errGetOrder)
	}
	order.ID = id

	if errGetDelivery := dataBase.conn.GetContext(ctx, &order.Delivery, `
	SELECT name, phone, zip, city, address, region, email
	FROM deliver
	WHERE order_id = $1
	`, order.ID); errGetDelivery != nil {
		return Models.OrderInfo{}, fmt.Errorf("get delivery: %v", errGetDelivery)
	}

	if errGetPayments := dataBase.conn.GetContext(ctx, &order.Payment, `
	SELECT transaction, request_id, currency, provider, amount,
                    payment_dt, bank, delivery_cost, goods_total, custom_fee
	FROM payment
	WHERE order_id = $1
	`, order.ID); errGetPayments != nil {
		return Models.OrderInfo{}, fmt.Errorf("get payments: %v", errGetPayments)
	}

	if errGetItems := dataBase.conn.SelectContext(ctx, &order.Items, `
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

	fmt.Println("start get lasts")
	if errGetOrder := dataBase.conn.SelectContext(ctx, &Orders, `
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
		if errGetDelivery := dataBase.conn.GetContext(ctx, &Orders[i].Delivery, `
			SELECT name, phone, zip, city, address, region, email
			FROM deliver
			WHERE order_id = $1`, Orders[i].ID);
			errGetDelivery != nil {
			fmt.Println("get delivery ", errGetDelivery)
			return nil, fmt.Errorf("get delivery: %v", errGetDelivery)
		}

		fmt.Println("delivery got")

		if errGetPayments := dataBase.conn.GetContext(ctx, &Orders[i].Payment, `
			SELECT transaction, request_id, currency, provider, amount,
                    payment_dt, bank, delivery_cost, goods_total, custom_fee
			FROM payment
			WHERE order_id = $1`, Orders[i].ID);
			errGetPayments != nil {
			fmt.Println("get payments ", errGetPayments)
			return nil, fmt.Errorf("get payments: %v", errGetPayments)
		}

		fmt.Println("payments got ", Orders[i].Payment)

		if errGetItems := dataBase.conn.SelectContext(ctx, &Orders[i].Items, `
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

func (dataBase *DB) InTx(ctx context.Context, isolation sql.IsolationLevel, f func(tx sqlx.Tx) error) error {
	tx, err := dataBase.conn.BeginTxx(ctx, &sql.TxOptions{
		Isolation: isolation,
	})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	if err := f(*tx); err != nil {
		if errRoll := tx.Rollback(); errRoll != nil {
			return fmt.Errorf("rollback tx: %v (error: %w)", errRoll, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
