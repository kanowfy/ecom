package main

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
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

type CartModel struct {
	DB *redis.Client
}

func NewModel(client *redis.Client) CartModel {
	return CartModel{
		DB: client,
	}
}

func (m CartModel) GetCart(userid string) (*Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	res, err := m.DB.HKeys(ctx, userid).Result()
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, ErrNoRecord
	}

	var items []*CartItem

	for _, pid := range res {
		q, err := m.DB.HGet(ctx, userid, pid).Result()
		if err != nil {
			return nil, err
		}

		qInt, err := strconv.ParseUint(q, 10, 32)
		if err != nil {
			return nil, err
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

func (m CartModel) AddItem(userid string, item *CartItem) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := m.GetCart(userid)
	if err != nil {
		return err
	}

	ok, err := m.DB.HSetNX(ctx, userid, item.ProductId, item.Quantity).Result()
	if err != nil {
		return err
	}

	if !ok {
		return ErrInCart
	}

	return nil
}

func (m CartModel) RemoveItem(userid, productid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := m.GetCart(userid)
	if err != nil {
		return err
	}

	return m.DB.HDel(ctx, userid, productid).Err()
}
