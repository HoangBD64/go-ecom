package interfaces

import (
	"context"

	"github.com/HoangBD64/go-ecom/pkg/api/handler/request"
	"github.com/HoangBD64/go-ecom/pkg/api/handler/response"
)

type StockUseCase interface {
	GetAllStockDetails(ctx context.Context, pagination request.Pagination) (stocks []response.Stock, err error)
	UpdateStockBySKU(ctx context.Context, updateDetails request.UpdateStock) error
}
