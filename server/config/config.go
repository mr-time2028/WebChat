package config

import "github.com/mr-time2028/WebChat/server/database"

type Config struct {
	Port   int
	DB     *database.DB
}
