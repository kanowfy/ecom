package repository

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

type Repository interface {
	GetById(ctx context.Context, id string) (*Product, error)
	Add(ctx context.Context, product Product) error
	List(ctx context.Context, query string) ([]Product, error)
}

type repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Repository {
	return &repository{
		db,
	}
}

func (r *repository) GetById(ctx context.Context, id string) (*Product, error) {
	stmt := `
	SELECT id, name, description, category, price, image
	FROM products
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	var product Product

	err := r.db.QueryRow(ctx, stmt, id).Scan(&product.ID, &product.Name, &product.Description, &product.Category, &product.Price, &product.Image)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecord
		default:
			return nil, fmt.Errorf("failed to query: %w", err)
		}
	}

	return &product, nil
}

func (r *repository) Add(ctx context.Context, product Product) error {
	stmt := `
	INSERT INTO products (id, name, description, category, price, image)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	args := []interface{}{product.ID, product.Name, product.Description, product.Category, product.Price, product.Image}

	if _, err := r.db.Exec(ctx, stmt, args...); err != nil {
		return fmt.Errorf("failed to add product: %w", err)
	}

	return nil
}

func (r *repository) List(ctx context.Context, query string) ([]Product, error) {
	document := "name || ' ' || description || ' ' || category"
	stmt := fmt.Sprintf(`
	SELECT id, name, description, category, price, image
	FROM products
	WHERE (to_tsvector('simple', %s) @@ plainto_tsquery($1) OR $1 = '')
	ORDER BY ts_rank(to_tsvector(%s), plainto_tsquery($1))
	`, document, document)

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	rows, err := r.db.Query(ctx, stmt, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	products := []Product{}

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Category, &product.Price, &product.Image)
		if err != nil {
			return nil, fmt.Errorf("error reading from row: %w", err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
