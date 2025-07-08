package api

import (
	"fmt"

	db "github.com/e-commerce/db/sqlc"
	"github.com/e-commerce/token"
	"github.com/e-commerce/util"
	"github.com/gin-gonic/gin"
)

// servers HTTP requests for the insta-app
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// Creates HTTP server and Setup Routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// Routes
	server.setupRoutes()

	return server, nil
}

func (server *Server) setupRoutes() {
	router := gin.Default()

	// routes
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	router.GET("/products", server.listProducts)
	router.GET("/products/:id", server.getProductByID)
	router.GET("/products/category/:id", server.listProductsByCategoryID)

	router.POST("/tokens/renew_access", server.renewAccessToken)

	// for both users and admins
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker,
		[]string{util.AdminRole, util.CustomerRole}))
	authRoutes.GET("/users/:id", server.getUserByID)
	authRoutes.DELETE("/users/:id", server.deleteUser)

	authRoutes.POST("/cart", server.addToCart)
	authRoutes.GET("/cart", server.getCartItemsByUser)
	authRoutes.DELETE("/cart/:id", server.deleteItemInCart)

	authRoutes.POST("/order-payment", server.createOrderAndPayment)
	authRoutes.GET("/order", server.listOrdersByUser)
	authRoutes.GET("/order/:id", server.getOrder)
	authRoutes.GET("/payment/:id", server.getPaymentByOrder)

	authRoutes.POST("/review", server.createReview)
	authRoutes.GET("/review/:id", server.getReviewsForProduct)

	// for only admins
	adminRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker,
		[]string{util.AdminRole}))
	adminRoutes.POST("/products", server.createProduct)
	adminRoutes.PATCH("/products/:id", server.updateProduct)
	adminRoutes.DELETE("/products/:id", server.deleteProduct)

	server.router = router

}

// Starts and runs HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
