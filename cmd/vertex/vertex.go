package main

import (
	"github.com/evernetlabs/vertex/internal/app/vertex"
	"github.com/evernetlabs/vertex/internal/pkg/env"
)

func main() {

	vertex.NewServer(&vertex.ServerConfig{
		Host:          env.GetOrDefault("HOST", "0.0.0.0"),
		Port:          env.GetOrDefault("PORT", "8080"),
		MongoEndpoint: env.GetOrDefault("MONGO_ENDPOINT", "mongodb://localhost:27017"),
		MongoDatabase: env.GetOrDefault("MONGO_DATABASE", "vertex"),
	}).Start()
}
