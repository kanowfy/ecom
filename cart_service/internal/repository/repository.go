package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ErrNoRecord = errors.New("no record")
	ErrInCart   = errors.New("item already in cart")
)

type CartItem struct {
	ProductId string
	Quantity  uint32
}

type Cart struct {
	UserId string
	Items  []*CartItem
}

type Repository interface {
	GetCart(userid string) (*Cart, error)
	RemoveItem(userid, productid string) error
	AddItem(userid string, item *CartItem) error
}

type repository struct {
	db *redis.Client
}

func New(client *redis.Client) Repository {
	return &repository{
		db: client,
	}
}

func (r *repository) GetCart(userid string) (*Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	res, err := r.db.HKeys(ctx, userid).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get cart result: %w", err)
	}

	if len(res) == 0 {
		return nil, ErrNoRecord
	}

	var items []*CartItem

	for _, pid := range res {
		q, err := r.db.HGet(ctx, userid, pid).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get cart item: %w", err)
		}

		qInt, err := strconv.ParseUint(q, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid quantity: %w", err)
		}

		items = append(items, &CartItem{
			ProductId: pid,
			Quantity:  uint32(qInt),
		})
	}

	return &Cart{
		UserId: userid,
		Items:  items,
	}, nil
}

func (r *repository) AddItem(userid string, item *CartItem) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := r.GetCart(userid)
	if err != nil {
		if !errors.Is(err, ErrNoRecord) {
			return fmt.Errorf("failed to get cart: %w", err)
		}
	}

	ok, err := r.db.HSetNX(ctx, userid, item.ProductId, item.Quantity).Result()
	if err != nil {
		return fmt.Errorf("failed to add item: %w", err)
	}

	if !ok {
		return ErrInCart
	}

	return nil
}

func (r *repository) RemoveItem(userid, productid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := r.GetCart(userid)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	if err = r.db.HDel(ctx, userid, productid).Err(); err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}
