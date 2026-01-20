package configs

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	}
	DB struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"db"`
	JWT JWTConfig `mapstructure:"jwt"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"signing_key"`
	AccessTTL  string `mapstructure:"access_ttl"`
	RefreshTTL string `mapstructure:"refresh_ttl"`
}
