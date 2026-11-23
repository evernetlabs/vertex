package vertex

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/evernetlabs/vertex/internal/app/vertex/account"
	"github.com/evernetlabs/vertex/internal/app/vertex/admin"
	"github.com/evernetlabs/vertex/internal/app/vertex/auth"
	"github.com/evernetlabs/vertex/internal/app/vertex/config"
	"github.com/evernetlabs/vertex/internal/app/vertex/health"
	"github.com/evernetlabs/vertex/internal/app/vertex/vertex"
	"github.com/evernetlabs/vertex/internal/pkg/exception"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Server struct {
	config *ServerConfig
}

func NewServer(config *ServerConfig) *Server {
	return &Server{config: config}
}

type ServerConfig struct {
	Host          string
	Port          string
	MongoEndpoint string
	MongoDatabase string
}

func (s *Server) Start() {
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(s.config.MongoEndpoint))
	if err != nil {
		log.Fatal("error connecting to mongodb: ", err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatal("error disconnecting from mongo: ", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("error connecting to mongodb: ", err)
	}

	mongoDatabase := mongoClient.Database(s.config.MongoDatabase)
	configsCollection := mongoDatabase.Collection("configs")
	adminsCollection := mongoDatabase.Collection("admins")
	accountsCollection := mongoDatabase.Collection("accounts")

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(exception.ErrorHandler())

	configManager := config.NewManager(configsCollection)
	authenticator := auth.NewAuthenticator(configManager)
	adminManager := admin.NewManager(adminsCollection, configManager, authenticator)
	accountManager := account.NewManager(accountsCollection, authenticator)

	health.NewHandler(router).Register()
	vertex.NewHandler(router, authenticator, configManager).Register()
	admin.NewHandler(router, authenticator, adminManager).Register()
	account.NewHandler(router, authenticator, accountManager).Register()

	err = router.Run(fmt.Sprintf("%s:%s", s.config.Host, s.config.Port))

	if err != nil {
		log.Fatalf("error starting server: %v\n", err)
	}
}
