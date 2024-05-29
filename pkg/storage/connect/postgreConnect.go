package connect

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

// NewPostgresConnection - функция создающая подключение к БД и предоставляющая его наружу
func NewPostgresConnection() (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig("postgresql://Ozon:ozon@pq_database:5432/Ozon?sslmode=disable")
	if err != nil {
		fmt.Println("Error parsing database source:", err)
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return nil, err
	}

	time.Sleep(3 * time.Second)

	err = pool.Ping(context.Background())
	if err != nil {
		fmt.Println("Error pinging the database:", err)
		return nil, err
	}

	fmt.Println("Db connection opened")
	return pool, nil
}
