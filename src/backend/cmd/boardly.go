package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/wuraLab/boardly/src/backend/internal/config"
	"github.com/wuraLab/boardly/src/backend/internal/controllers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	var DB *gorm.DB
	var err error
	var Config config.Config
	//Load configuration
	if Config, err = config.LoadConfig(".env"); err != nil {
		log.Fatal(err)
	}
	//connect the DB
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s  sslmode=disable", Config.Database.Host, Config.Database.Port, Config.Database.User, Config.Database.Password, Config.Database.DBName)
	if DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		log.Fatalln(err)
	}

	//Start the default gin server
	r := gin.Default()
	api := r.Group("/api")
	{
		/*** START USER ***/
		user := controllers.User{
			Base: controllers.Base{
				DB: DB,
			},
		}

		// v1.POST("/user/login", user.Login)
		api.POST("/user/register", user.Register)
		// v1.GET("/user/logout", user.Logout)
	}
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "Welcome to Boardly",
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{})
	})

	log.Infoln(Config.Server.SSL)

	if Config.Server.ENV == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}

	if Config.Server.SSL == "TRUE" {

		SSLKeys := &struct {
			CERT string
			KEY  string
		}{}

		//Generated using sh generate-certificate.sh
		SSLKeys.CERT = "./cert/myCA.cer"
		SSLKeys.KEY = "./cert/myCA.key"

		log.Fatal(r.RunTLS(":"+Config.Server.Port, SSLKeys.CERT, SSLKeys.KEY))
	} else {
		log.Fatal(r.Run(":" + Config.Server.Port))
	}
}
