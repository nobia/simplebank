package api

import (
	"fmt"

	db "github.com/nobia/simplebank/db/sqlc"
	"github.com/nobia/simplebank/token"
	"github.com/nobia/simplebank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store        db.Store
	config       util.Config
	tokenManager token.TokenManager
	router       *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenManager, err := token.NewPasetoManager(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token manager: %w", err)
	}
	server := &Server{
		store:        store,
		config:       config,
		tokenManager: tokenManager,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenManager))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)

	authRoutes.POST("/transfer", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
