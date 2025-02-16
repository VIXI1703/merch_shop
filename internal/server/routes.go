package server

import (
	"merch_shop/internal/handlers"
	"merch_shop/internal/middleware"
	"merch_shop/internal/provider"
	"merch_shop/internal/repository"
	"merch_shop/internal/service"
)

func (server *Server) ConfigureRoutes() {
	uow := repository.NewGormUnitOfWork(server.DB)
	jwtAuth := provider.NewJWTAuth([]byte(server.Cfg.JWT.SigningKey), server.Cfg.JWT.Duration)
	jwtMiddleware := middleware.JWTAuthMiddleware(jwtAuth)
	transactionService := service.NewTransactionService(uow)
	authService := service.NewAuthService(jwtAuth, uow)

	transactionHandler := handlers.NewTransactionHandler(transactionService)
	authHandler := handlers.NewAuthHandler(authService)
	apiRoute := server.Gin.Group("/api")
	authHandler.Routes(apiRoute)

	protectedRoutes := apiRoute.Group("/", jwtMiddleware)

	transactionHandler.Routes(protectedRoutes)

}
