package order

import (
	"database/sql"

	"github.com/gonan98/go-ecommerce-api/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateOrder(order types.Order) (int, error) {
	res, err := s.db.Exec("INSERT INTO orders (status, user_id) VALUES (?, ?)", order.Status, order.UserID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *Store) CreateOrderDetail(detail types.OrderDetail) error {
	_, err := s.db.Exec("INSERT INTO order_details (order_id, product_id, quantity, unit_price, discount) VALUES (?, ?, ?, ?, ?)", detail.OrderID, detail.ProductID, detail.Quantity, detail.UnitPrice, detail.Discount)
	return err
}
