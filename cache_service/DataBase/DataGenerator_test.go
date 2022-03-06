package DataBase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestConnectToDB(t *testing.T) {
	dataBase, err := sqlx.Connect("pgx", "postgresql://maui:maui@192.168.0.12:5432/postgres")
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < 500; i++ {
		var deliverID, paymentsID string
		var itemTest []item
		for j := 0; j < rand.Intn(5); j++{
			it := item {
				ChrtID: rand.Intn(10000),
				TrackNumber: "test " + strconv.Itoa(rand.Intn(10000)),
				Price: rand.Intn(10000),
				RID: "test " + strconv.Itoa(rand.Intn(10000)),
				Name: "test " + strconv.Itoa(rand.Intn(10000)),
				Sale: rand.Intn(10000),
				Size: "test " + strconv.Itoa(rand.Intn(10000)),
				TotalPrice: rand.Intn(10000),
				NmID: rand.Intn(10000),
				Brand: "test " + strconv.Itoa(rand.Intn(10000)),
				Status: rand.Intn(10000),
			}

			itemTest = append(itemTest, it)
		}

		order := orderI{
			OrderUID: "test " + strconv.Itoa(rand.Intn(10000)),
			TrackNumber: "test " + strconv.Itoa(rand.Intn(10000)),
			Entry: "test " + strconv.Itoa(rand.Intn(10000)),
			Delivery: deliv{
				Name: "test " + strconv.Itoa(rand.Intn(10000)),
				Phone: "test " + strconv.Itoa(rand.Intn(10000)),
				Zip: "test " + strconv.Itoa(rand.Intn(10000)),
				City: "test " + strconv.Itoa(rand.Intn(10000)),
				Address: "test " + strconv.Itoa(rand.Intn(10000)),
				Region: "test " + strconv.Itoa(rand.Intn(10000)),
				Email: "test " + strconv.Itoa(rand.Intn(10000)),
			},
			Payment: payment{
				Transaction: "test " + strconv.Itoa(rand.Intn(10000)),
				RequestID: "test " + strconv.Itoa(rand.Intn(10000)),
				Currency: "test " + strconv.Itoa(rand.Intn(10000)),
				Provider: "test " + strconv.Itoa(rand.Intn(10000)),
				Amount: rand.Intn(10000),
				PaymentDt: rand.Intn(10000),
				Bank: "test " + strconv.Itoa(rand.Intn(10000)),
				DeliveryCost: rand.Intn(10000),
				GoodsTotal: rand.Intn(10000),
				CustomFee: rand.Intn(10000),
			},
			Items: itemTest,
			Locale: "test " + strconv.Itoa(rand.Intn(10000)),
			InternalSignature: "test " + strconv.Itoa(rand.Intn(10000)),
			CustomerID: "test " + strconv.Itoa(rand.Intn(10000)),
			DeliveryService: "test " + strconv.Itoa(rand.Intn(10000)),
			ShardKey: "test " + strconv.Itoa(rand.Intn(10000)),
			SmID: rand.Intn(10000),
			DateCreated: time.Now(),
			OofShard: "test " + strconv.Itoa(rand.Intn(10000)),
		}

		if errTx := InTx(context.Background(), dataBase, sql.LevelReadCommitted, func(tx sqlx.Tx) error {
			fmt.Println("start tx")

			var orderID string
			errExecOrder := tx.Get(&orderID, `
				INSERT INTO order_info (order_uid, track_number, entry, locale, internal_signature,
                     customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
				RETURNING id`, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
				order.CustomerID, order.DeliveryService, order.ShardKey,
				order.SmID, order.DateCreated, order.OofShard)
			if errExecOrder != nil {
				fmt.Println("failed order insert: ", errExecOrder)
			}
			order.ID = orderID

			if errExecDelivery := tx.Get(&deliverID, `
				INSERT INTO deliver(order_id, name, phone, zip, city, address, region, email)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING id;`, order.ID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City,
				order.Delivery.Address, order.Delivery.Region, order.Delivery.Email);
			errExecDelivery != nil {
				fmt.Println("failed delivery insert: ", errExecDelivery)
			}
			fmt.Println("delivery end")

			if errExecPayments := tx.Get(&paymentsID, `
				INSERT INTO payment(order_id, transaction, request_id, currency, provider, amount,
                 payment_dt, bank, delivery_cost, goods_total, custom_fee)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
				RETURNING id;`, order.ID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
				order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt,
				order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee);
			errExecPayments != nil {
				fmt.Println("failed payment insert: ", errExecPayments)
			}
			fmt.Println("payment end")

			for _, v := range order.Items {
				_, errExecItem := tx.Exec(`
					INSERT INTO item(order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`, order.ID, v.ChrtID, v.TrackNumber,
					v.Price, v.RID, v.Name, v.Sale, v.Size, v.TotalPrice, v.NmID, v.Brand, v.Status)
				if errExecItem != nil {
					fmt.Println("failed insert item: ", errExecItem)
				}
			}
			fmt.Println("item end")

			return nil
		}); errTx != nil {
			fmt.Println("errTx")
		}
	}

	fmt.Println("tx complete")
}

func InTx(ctx context.Context, dataBase *sqlx.DB, isolation sql.IsolationLevel, f func(tx sqlx.Tx) error) error {
	tx, err := dataBase.BeginTxx(ctx, &sql.TxOptions{
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

type orderI struct {
	ID                string    `json:"id" db:"id"`
	OrderUID          string    `json:"order_uid,omitempty" db:"order_uid"`
	TrackNumber       string    `json:"track_number,omitempty" db:"track_number"`
	Entry             string    `json:"entry,omitempty" db:"entry"`
	Delivery          deliv  `json:"delivery"`
	Payment           payment   `json:"payment"`
	Items             []item    `json:"items,omitempty"`
	Locale            string    `json:"locale,omitempty" db:"locale"`
	InternalSignature string    `json:"internal_signature,omitempty" db:"internal_signature"`
	CustomerID        string    `json:"customer_id,omitempty" db:"customer_id"`
	DeliveryService   string    `json:"delivery_service,omitempty" db:"delivery_service"`
	ShardKey          string    `json:"shard_key,omitempty" db:"shard_key"`
	SmID              int       `json:"sm_id,omitempty" db:"sm_id"`
	DateCreated       time.Time `json:"date_created" db:"date_created"`
	OofShard          string    `json:"oof_shard,omitempty" db:"oof_shard"`
}

type deliv struct {
	Name    string `json:"name,omitempty" db:"name"`
	Phone   string `json:"phone,omitempty" db:"phone"`
	Zip     string `json:"zip,omitempty" db:"zip"`
	City    string `json:"city,omitempty" db:"city"`
	Address string `json:"address,omitempty" db:"address"`
	Region  string `json:"region,omitempty" db:"region"`
	Email   string `json:"email,omitempty" db:"email"`
}

type payment struct {
	Transaction  string `json:"transaction,omitempty" db:"transaction"`
	RequestID    string `json:"request_id,omitempty" db:"request_id"`
	Currency     string `json:"currency,omitempty" db:"currency"`
	Provider     string `json:"provider,omitempty" db:"provider"`
	Amount       int    `json:"amount,omitempty" db:"amount"`
	PaymentDt    int    `json:"payment_dt,omitempty" db:"payment_dt"`
	Bank         string `json:"bank,omitempty" db:"bank"`
	DeliveryCost int    `json:"delivery_cost,omitempty" db:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total,omitempty" db:"goods_total"`
	CustomFee    int    `json:"custom_fee,omitempty" db:"custom_fee"`
}

type item struct {
	ChrtID      int    `json:"chrt_id,omitempty" db:"chrt_id"`
	TrackNumber string `json:"track_number,omitempty" db:"track_number"`
	Price       int    `json:"price,omitempty" db:"price"`
	RID         string `json:"rid,omitempty" db:"rid"`
	Name        string `json:"name,omitempty" db:"name"`
	Sale        int    `json:"sale,omitempty" db:"sale"`
	Size        string `json:"size,omitempty" db:"size"`
	TotalPrice  int    `json:"total_price,omitempty" db:"total_price"`
	NmID        int    `json:"nm_id,omitempty" db:"nm_id"`
	Brand       string `json:"brand,omitempty" db:"brand"`
	Status      int    `json:"status,omitempty" db:"status"`
}
