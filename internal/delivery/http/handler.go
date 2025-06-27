package http

import (
	"github.com/elokanugrah/go-order-system/internal/usecase"
)

type Handler struct {
	productUseCase *usecase.ProductUseCase
	orderUseCase   *usecase.OrderUseCase
}

func NewHandler(puc *usecase.ProductUseCase, ouc *usecase.OrderUseCase) *Handler {
	return &Handler{
		productUseCase: puc,
		orderUseCase:   ouc,
	}
}
