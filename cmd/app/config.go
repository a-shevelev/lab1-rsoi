package main

import (
	"lab1-rsoi/internal/server"
	"lab1-rsoi/pkg/postgres"
)

type Config struct {
	Server server.Server   `envconfig:"SERVER"`
	DB     postgres.Config `envconfig:"DB"`
}
