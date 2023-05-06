package rest

import (
	"food_delivery_api/pkg/middleware"
	"food_delivery_api/pkg/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func Handler(s service.Service) *gin.Engine {
	r := gin.Default()
	_ = r.SetTrustedProxies(nil)
	setupCORS(r)

	// Public API
	r.GET("/health", getHealthStatus)
	r.POST("/api/v1/login", login(s))

	// Protected API
	v1 := r.Group("/api/v1")
	v1.Use(middleware.JWT())
	{
		// Users
		v1.POST("/users", addUser(s))
		v1.GET("/users", getUsers(s))
		v1.GET("/users/:id", getUser(s))
		v1.GET("/users/me", getLoggedUser(s))
		v1.PUT("/users/:id", editUser(s))
		v1.DELETE("/users/:id", removeUser(s))

		// Categories
		v1.POST("/categories", addCategory(s))
		v1.GET("/categories", getCategories(s))
		v1.GET("/categories/:id", getCategory(s))
		v1.PUT("/categories/:id", editCategory(s))
		v1.DELETE("/categories/:id", removeCategory(s))

		// UOMs
		v1.POST("/uoms", addUOM(s))
		v1.GET("/uoms", getUOMs(s))
		v1.GET("/uoms/:id", getUOM(s))
		v1.PUT("/uoms/:id", editUOM(s))
		v1.DELETE("/uoms/:id", removeUOM(s))

		// Products
		v1.POST("/products", addProduct(s))
		v1.GET("/products", getProducts(s))
		v1.GET("/products/:id", getProduct(s))
		v1.PUT("/products/:id", editProduct(s))
		v1.DELETE("/products/:id", removeProduct(s))

		// Transactions
		v1.POST("/transactions", addTransaction(s))
		v1.GET("/transactions", getTransactions(s))
		v1.GET("/transactions/:id", getTransaction(s))
		v1.PUT("/transactions/:id", editTransaction(s))
		v1.DELETE("/transactions/:id", removeTransaction(s))
	}

	return r
}

func setupCORS(r *gin.Engine) {
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, PATCH, POST, DELETE, OPTIONS",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          1 * time.Minute,
		Credentials:     false,
		ValidateHeaders: false,
	}))
}

func getHealthStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
