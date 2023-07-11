package handlers

import (
	"food_delivery_api/cfg"
	"food_delivery_api/pkg/model"
	"food_delivery_api/pkg/service"
	"food_delivery_api/pkg/storage/mysql"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func main() {
	goEnv := strings.ToLower(os.Getenv("GO_ENV"))
	if goEnv == "" {
		goEnv = "local"
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Load the config
	LoadConfig(goEnv)

	log.Print(strings.ToUpper(cfg.AppName), " is warming up ...")

	// Run the server
	// MySQL setup
	rmy, err := mysql.NewStorage(cfg.My, goEnv)
	if err != nil {
		log.Fatal("Error: Database failed to connect (", cfg.My.DSN, ") - ", err)
	}

	// Handler setup
	s := service.NewService(rmy)

	r := gin.Default()
	_ = r.SetTrustedProxies(nil)
	setupCORS(r)

	// Public API
	// r.GET("/health", getHealthStatus)
	// r.POST("/api/v1/login", rest.Login(s))

	// Protected API
	v1 := r.Group("/api/v1")
	// v1.Use(middleware.JWT())
	{
		// Users
		// v1.POST("/users", rest.AddUser(s))
		// v1.POST("/users/upload", rest.AddUsers(s))
		v1.GET("/users", GetUsers(s))
		// v1.GET("/users/:id", rest.GetUser(s))
		// v1.GET("/users/me", rest.GetLoggedUser(s))
		// v1.PUT("/users/:id", rest.EditUser(s))
		// v1.DELETE("/users/:id", rest.RemoveUser(s))

		// // Categories
		// v1.POST("/categories", addCategory(s))
		// v1.POST("/categories/upload", addCategories(s))
		// v1.GET("/categories", getCategories(s))
		// v1.GET("/categories/:id", getCategory(s))
		// v1.PUT("/categories/:id", editCategory(s))
		// v1.DELETE("/categories/:id", removeCategory(s))

		// // UOMs
		// v1.POST("/uoms", addUOM(s))
		// v1.POST("/uoms/upload", addUOMs(s))
		// v1.GET("/uoms", getUOMs(s))
		// v1.GET("/uoms/:id", getUOM(s))
		// v1.PUT("/uoms/:id", editUOM(s))
		// v1.DELETE("/uoms/:id", removeUOM(s))

		// // Products
		// v1.POST("/products", addProduct(s))
		// v1.POST("/products/upload", addProducts(s))
		// v1.GET("/products", getProducts(s))
		// v1.GET("/products/:id", getProduct(s))
		// v1.PUT("/products/:id", editProduct(s))
		// v1.DELETE("/products/:id", removeProduct(s))

		// // Transactions
		// v1.POST("/transactions", addTransaction(s))
		// v1.GET("/transactions", getTransactions(s))
		// v1.GET("/transactions/:id", getTransaction(s))
		// v1.PUT("/transactions/:id", editTransaction(s))
		// v1.DELETE("/transactions/:id", removeTransaction(s))

		// // Reports
		// v1.GET("/reports/dashboard", getReportDashboard(s))
		// v1.GET("/reports/chart", getReportChart(s))
		// v1.GET("/reports/exponential-smoothing", getReportExponentialSmoothing(s))
		// v1.GET("/reports/monthly-exponential-smoothing", getReportMonthlyExponentialSmoothing(s))
	}

	host := cfg.Glb.Serv.Host
	if host == "" {
		host = GetLocalIP()
		cfg.Glb.Serv.Host = host
	}

	if goEnv == "local" {
		log.Println("Server Running on", goEnv, "environment, (REST APIs) listening on", host+":"+cfg.Serv.Port)
		log.Fatal("Error: Server failed to run - ", r.Run(cfg.Serv.Host+":"+cfg.Serv.Port))
	} else {
		log.Println("Server Running on", goEnv, "environment, (REST APIs) listening on", host+":"+os.Getenv("PORT"))
		log.Fatal("Error: Server failed to run - ", r.Run(":"+os.Getenv("PORT")))
	}
}

func LoadConfig(goEnv string) {
	var arg string

	if config := os.Getenv("CONFIG_FILE"); config != "" {
		arg = config
	} else if len(os.Args) == 2 {
		arg = "cfg/config." + os.Args[1] + ".yaml"
	} else {
		arg = "cfg/config." + goEnv + ".yaml"
	}

	err := cfg.Load(arg)
	if err != nil {
		log.Fatal("Error: config failed to load - ", err)
	}

	log.Println("Load config from", arg)
}

func GetLocalIP() string {
	address, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range address {
		if inet, ok := address.(*net.IPNet); ok && !inet.IP.IsLoopback() {
			if inet.IP.To4() != nil {
				return inet.IP.String()
			}
		}
	}

	return ""
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

func GetUsers(s service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var res []model.User
		// var ttl int64
		var err error

		qp := model.QueryPagination{}
		err = c.ShouldBindQuery(&qp)
		if err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		if res, _, err = s.GetUsers(qp); err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}
