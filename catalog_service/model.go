package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNoRecord = errors.New("no record")
)

type Product struct {
	ID          string
	Name        string
	Description string
	Category    string
	Price       int64
	Image       string
}

type ProductModel struct {
	DB *pgxpool.Pool
}

func NewModel(db *pgxpool.Pool) ProductModel {
	return ProductModel{
		DB: db,
	}
}

func (m ProductModel) GetById(id string) (*Product, error) {
	stmt := `
	SELECT id, name, description, category, price, image
	FROM products
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var product Product

	err := m.DB.QueryRow(ctx, stmt, id).Scan(&product.ID, &product.Name, &product.Description, &product.Category, &product.Price, &product.Image)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecord
		default:
			return nil, err
		}
	}

	return &product, nil
}

func (m ProductModel) Add(product *Product) error {
	stmt := `
	INSERT INTO products (id, name, description, category, price, image)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	args := []interface{}{product.ID, product.Name, product.Description, product.Category, product.Price, product.Image}

	_, err := m.DB.Exec(ctx, stmt, args...)
	return err
}

func (m ProductModel) List(query string) ([]*Product, error) {
	document := "name || ' ' || description || ' ' || category"
	stmt := fmt.Sprintf(`
	SELECT id, name, description, category, price, image
	FROM products
	WHERE (to_tsvector('simple', %s) @@ plainto_tsquery($1) OR $1 = '')
	ORDER BY ts_rank(to_tsvector(%s), plainto_tsquery($1))
	`, document, document)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := m.DB.Query(ctx, stmt, query)
	if err != nil {
		return nil, err
	}

	products := []*Product{}

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Category, &product.Price, &product.Image)
		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
