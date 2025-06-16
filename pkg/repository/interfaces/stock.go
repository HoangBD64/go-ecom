package interfaces

import (
	"context"

	"github.com/HoangBD64/go-ecom/pkg/api/handler/request"
	"github.com/HoangBD64/go-ecom/pkg/api/handler/response"
)

type StockRepository interface {
	FindAll(ctx context.Context, pagination request.Pagination) (stocks []response.Stock, err error)
	Update(ctx context.Context, updateValues request.UpdateStock) error
}
