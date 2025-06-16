//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	http "github.com/HoangBD64/go-ecom/pkg/api"
	"github.com/HoangBD64/go-ecom/pkg/api/handler"
	"github.com/HoangBD64/go-ecom/pkg/api/middleware"
	"github.com/HoangBD64/go-ecom/pkg/config"
	"github.com/HoangBD64/go-ecom/pkg/db"
	"github.com/HoangBD64/go-ecom/pkg/repository"
	"github.com/HoangBD64/go-ecom/pkg/service/cloud"
	"github.com/HoangBD64/go-ecom/pkg/service/otp"
	"github.com/HoangBD64/go-ecom/pkg/service/token"
	"github.com/HoangBD64/go-ecom/pkg/usecase"
)

func InitializeApi(cfg config.Config) (*http.ServerHTTP, error) {

	wire.Build(db.ConnectDatabase,
		//external
		token.NewTokenService,
		otp.NewOtpAuth,
		cloud.NewAWSCloudService,

		// repository

		middleware.NewMiddleware,
		repository.NewAuthRepository,
		repository.NewPaymentRepository,
		repository.NewAdminRepository,
		repository.NewUserRepository,
		repository.NewCartRepository,
		repository.NewProductRepository,
		repository.NewOrderRepository,
		repository.NewCouponRepository,
		repository.NewOfferRepository,
		repository.NewStockRepository,
		repository.NewBrandDatabaseRepository,

		//usecase
		usecase.NewAuthUseCase,
		usecase.NewAdminUseCase,
		usecase.NewUserUseCase,
		usecase.NewCartUseCase,
		usecase.NewPaymentUseCase,
		usecase.NewProductUseCase,
		usecase.NewOrderUseCase,
		usecase.NewCouponUseCase,
		usecase.NewOfferUseCase,
		usecase.NewStockUseCase,
		usecase.NewBrandUseCase,
		// handler
		handler.NewAuthHandler,
		handler.NewAdminHandler,
		handler.NewUserHandler,
		handler.NewCartHandler,
		handler.NewPaymentHandler,
		handler.NewProductHandler,
		handler.NewOrderHandler,
		handler.NewCouponHandler,
		handler.NewOfferHandler,
		handler.NewStockHandler,
		handler.NewBrandHandler,

		http.NewServerHTTP,
	)

	return &http.ServerHTTP{}, nil
}
