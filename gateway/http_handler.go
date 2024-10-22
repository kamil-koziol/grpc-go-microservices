package main

import (
	"errors"
	"net/http"

	"github.com/kamil-koziol/common"
	pb "github.com/kamil-koziol/common/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	client pb.OrderServiceClient
}

func NewHandler(client pb.OrderServiceClient) *handler {
	return &handler{client}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerId}/orders", h.HandleCreateOrder)
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerId := r.PathValue("customerId")

	var items []*pb.ItemsWithQuantity
	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateItems(items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	o, err := h.client.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerId: customerId,
		Items:      items,
	})
	rStatus := status.Convert(err)
	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}
	}

	common.WriteJSON(w, http.StatusOK, o)
}

func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return errors.New("items must have at least one item")
	}

	for _, i := range items {
		if i.ID == "" {
			return errors.New("Item id is required")
		}

		if i.Quantity <= 0 {
			return errors.New("items must have a valid quantity")
		}
	}

	return nil
}
