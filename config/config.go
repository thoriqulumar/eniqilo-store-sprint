package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	DB         DBConfig `env:",prefix=DB_,required"`
	BcryptSalt int      `env:"BCRYPT_SALT"`
	JWTSecret  string   `env:"JWT_SECRET"`
}

type DBConfig struct {
	Name     string `env:"NAME"`
	Port     string `env:"PORT"`
	Host     string `env:"HOST"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Params   string `env:"PARAMS"`
}

func LoadConfig(ctx context.Context) (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c DBConfig) ConnectionString() string {
	params := strings.ReplaceAll(c.Params, `"`, "")
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?%s", c.Username, c.Password, c.Host, c.Port, c.Name, params)
}
