package db

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Options struct {
	DSN         string
	MaxOpen     int
	MaxIdle     int
	MaxIdleTime time.Duration
	// Лог-уровень для dev/prod
	LogLevel logger.LogLevel // logger.Info на dev, logger.Warn на prod
	// Префиксы/сингуларизация имён таблиц — опционально
	TablePrefix   string
	SingularTable bool
}

type Conn struct {
	Gorm *gorm.DB
	SQL  *sql.DB
}

func Open(opts Options) (*Conn, error) {
	gdb, err := gorm.Open(postgres.Open(opts.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(opts.LogLevel),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   opts.TablePrefix,
			SingularTable: opts.SingularTable,
		},
	})
	if err != nil {
		return nil, err
	}

	sdb, err := gdb.DB()
	if err != nil {
		return nil, err
	}

	if opts.MaxOpen > 0 {
		sdb.SetMaxOpenConns(opts.MaxOpen)
	}
	if opts.MaxIdle > 0 {
		sdb.SetMaxIdleConns(opts.MaxIdle)
	}
	if opts.MaxIdleTime > 0 {
		sdb.SetConnMaxIdleTime(opts.MaxIdleTime)
	}

	// проверочный пинг
	if err := sdb.Ping(); err != nil {
		return nil, err
	}
	return &Conn{Gorm: gdb, SQL: sdb}, nil
}

func (c *Conn) HealthPing(ctx context.Context) error {
	return c.SQL.PingContext(ctx)
}
