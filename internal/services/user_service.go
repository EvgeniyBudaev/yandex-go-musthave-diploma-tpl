package services

import (
	"context"
	"fmt"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/storage"
	"strconv"
)

type UserService struct {
	db storage.Storage
}

func NewUserService(storage storage.Storage) *UserService {
	return &UserService{db: storage}
}

func (us *UserService) GetOrdersNotProcessed() ([]string, error) {
	orders, err := us.db.GetOrdersInProgress(context.Background())
	if err != nil {
		return nil, err
	}
	numbers := make([]string, 0)
	for _, order := range orders {
		num, err := strconv.ParseInt(order.Number, 10, 64)
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, strconv.FormatInt(num, 10))
	}
	if len(numbers) == 0 {
		return nil, fmt.Errorf("empty list")
	}
	return numbers, nil
}

func (us *UserService) UpdateOrder(accrual *storage.AccrualDto) error {
	order := storage.AccrualDto{Order: accrual.Order, Accrual: accrual.Accrual, Status: accrual.Status}
	err := us.db.UpdateOrder(context.Background(), order)
	if err != nil {
		return err
	}
	return nil
}
