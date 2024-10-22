package main

import (
	"context"
	"log"

	"github.com/kamil-koziol/common"
	pb "github.com/kamil-koziol/common/api"
)

type service struct {
	store OrdersStore
}

func NewService(store OrdersStore) *service {
	return &service{
		store: store,
	}
}

func (s *service) CreateOrder(ctx context.Context) error {
	return nil
}

func (s *service) ValidateOrder(ctx context.Context, payload *pb.CreateOrderRequest) error {
	if len(payload.Items) == 0 {
		return common.ErrNoItems
	}

	mergedItems := mergeItemsQuantities(payload.Items)
	log.Println(mergedItems)

	// validate with the stock service
	return nil
}

func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity, 0)

	counter := make(map[string]int32)
	for _, i := range items {
		counter[i.ID]++
	}

	for _, i := range items {
		itemCount, exists := counter[i.ID]
		if !exists {
			continue
		}

		i.Quantity = itemCount
		delete(counter, i.ID)

		merged = append(merged, i)
	}

	return merged
}
