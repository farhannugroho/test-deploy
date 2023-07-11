package main

import (
	"food_delivery_api/cfg"
	"food_delivery_api/pkg/http/rest"
	"food_delivery_api/pkg/service"
	"food_delivery_api/pkg/storage/mysql"
	"log"
	"net"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Handler() {
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
	run(goEnv)
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

func run(goEnv string) {
	// MySQL setup
	rmy, err := mysql.NewStorage(cfg.My, goEnv)
	if err != nil {
		log.Fatal("Error: Database failed to connect (", cfg.My.DSN, ") - ", err)
	}

	// Handler setup
	s := service.NewService(rmy)
	r := rest.Handler(s)

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

func main() {
	Handler()
}