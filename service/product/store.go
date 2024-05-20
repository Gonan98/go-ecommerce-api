package product

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gonan98/go-ecommerce-api/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetAll() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	products := make([]types.Product, 0)
	for rows.Next() {
		p, err := rowsToProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil
}

func (s *Store) GetByID(ID int) (*types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products WHERE id=?", ID)
	if err != nil {
		return nil, err
	}

	product := new(types.Product)
	for rows.Next() {
		product, err = rowsToProduct(rows)
		if err != nil {
			return nil, err
		}
	}

	if product.ID == 0 {
		return nil, fmt.Errorf("product %d not found", ID)
	}

	return product, nil
}

func (s *Store) GetByIDs(IDs []int) ([]types.Product, error) {
	placeholders := strings.Repeat(",?", len(IDs)-1)
	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (?%s)", placeholders)

	args := make([]interface{}, len(IDs))
	for i, v := range IDs {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	products := []types.Product{}
	for rows.Next() {
		p, err := rowsToProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil
}

func (s *Store) Create(product types.Product) error {
	_, err := s.db.Exec("INSERT INTO products (name, description, brand, unit_price, stock) VALUES (?,?,?,?,?)", product.Name, product.Description, product.Brand, product.UnitPrice, product.Stock)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Update(product types.Product) error {
	_, err := s.db.Exec("UPDATE products SET name=?, description=?, brand=?, unit_price=?, stock=? WHERE id=?", product.Name, product.Description, product.Brand, product.UnitPrice, product.Stock, product.ID)
	if err != nil {
		return err
	}

	return nil
}

func rowsToProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)
	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Brand,
		&product.UnitPrice,
		&product.Stock,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}
