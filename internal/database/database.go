package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Romasmi/go-rest-api-template/internal/config"
	"github.com/Romasmi/go-rest-api-template/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbConnection struct {
	DB     *pgxpool.Pool
	Config *config.Config
}

func (c *DbConnection) Connect() error {
	pgConfig, err := pgxpool.ParseConfig(c.Config.Database.URL)
	if err != nil {
		return fmt.Errorf("unable to parse database URL: %w", err)
	}
	pgConfig.MaxConns = int32(c.Config.Database.MaxConnections)
	pgConfig.MinConns = int32(c.Config.Database.MinConnections)
	pgConfig.MaxConnLifetime = utils.MinutesToNanoseconds(c.Config.Database.MaxConnectionLifetime)
	pgConfig.MaxConnIdleTime = utils.MinutesToNanoseconds(c.Config.Database.MaxConnectionIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c.DB, err = pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := c.Ping(); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	return nil
}

func (c *DbConnection) Close() {
	if c.DB != nil {
		c.DB.Close()
		log.Println("database connection closed")
	}
}

func (c *DbConnection) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return c.DB.Ping(ctx)
}
