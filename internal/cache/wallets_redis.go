package cache

import "time"

type Config struct {
	Addr         string    // Redis server address (host:port)
	Password     string    // No password set
	DB           string    // Use the default DB
	PoolSize     int       // Connection pool size
	DialTimeout  time.Time // Timeout for establishing the connection
	ReadTimeout  time.Time // Timeout for read operations
	WriteTimeout time.Time // Timeout for write operations
}
