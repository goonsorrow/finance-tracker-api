package configs

type Config struct {
	Server Server    `mapstructre:"server"`
	DB     DB        `mapstructure:"db"`
	Redis  Redis     `mapstructure:"redis"`
	JWT    JWTConfig `mapstructure:"jwt"`
}

type Server struct {
	Port string
}

type DB struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}
type JWTConfig struct {
	SigningKey string `mapstructure:"signing_key"`
	AccessTTL  string `mapstructure:"access_ttl"`
	RefreshTTL string `mapstructure:"refresh_ttl"`
}
