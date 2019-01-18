package mssql

import (
	"bytes"
	"fmt"
	"gitlab/nefco/platform/codegen/helper"
	"gitlab/nefco/platform/codegen/schema"
	"go.uber.org/zap"
	"strings"

	"github.com/99designs/gqlgen/handler"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gobuffalo/packr"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
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

func (s *mssql) connection(v *viper.Viper) (*sqlx.DB, error) {
	cfg := DefaultConfig
	if err := v.UnmarshalKey("app.service.mssql", &cfg); err != nil {
		return nil, err
	}

	logger := zap.L().Named("mssql")
	logger.Info("mssql config", zap.Any("cfg", cfg))

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
	return db, nil
}

func (s *mssql) Middleware(v *viper.Viper) (handler.Option, error) {
	db, err := s.connection(v)
	if err != nil {
		return nil, err
	}
	return handler.RequestMiddleware(middleware(db)), nil
}

func (s *mssql) Init(schemaPath string,v *viper.Viper) (error) {
	tmpl, err := helper.ReadTemplate("migration/migration", s.Box())
	if err != nil {
		return err
	}

	sch, err := schema.LoadSchema(schemaPath)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer([]byte{})
	if err := tmpl.Execute(buf, sch); err != nil {
		return err
	}
	query := buf.String()

	if !v.GetBool("prod") {
		outputMigrate := v.GetString("codegen.output.dir") + "migration.sql"
		if err := helper.WriteFile(outputMigrate, buf); err != nil {
			return err
		}
	}

	db, err := s.connection(v)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(query)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
