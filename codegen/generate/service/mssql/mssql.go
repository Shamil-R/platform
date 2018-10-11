package mssql

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"

	"github.com/99designs/gqlgen/handler"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gobuffalo/packr"
	"github.com/spf13/viper"
)

type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

var DefaultConfig = Config{
	Host:     "127.0.0.1",
	Port:     1433,
	Username: "username",
	Password: "password",
	Database: "database",
}

type mssql struct{}

func New() *mssql {
	return &mssql{}
}

func (s *mssql) Name() string {
	return "mssql"
}

func (s *mssql) Box() packr.Box {
	return packr.NewBox("./templates")
}

func (s *mssql) Middleware(v *viper.Viper) (handler.Option, error) {
	cfg := DefaultConfig
	if err := v.UnmarshalKey("app.service.mssql", &cfg); err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
	db, err := sqlx.Connect("mssql", dsn)
	if err != nil {
		return nil, err
	}
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	return handler.RequestMiddleware(middleware(db)), nil
}
