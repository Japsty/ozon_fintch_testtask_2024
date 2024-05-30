package connect

import (
	"database/sql"
	"fmt"
	"time"
)

// NewPostgresConnection - функция создающая подключение к БД и предоставляющая его наружу
func NewPostgresConnection() (*sql.DB, error) {
	conn, err := sql.Open("postgres", "postgresql://Ozon:ozon@pq_database:5432/Ozon?sslmode=disable")
	if err != nil {
		fmt.Println("Error parsing database config", err)
		return nil, err
	}

	time.Sleep(3 * time.Second)

	err = conn.Ping()
	if err != nil {
		fmt.Println("Error pinging the database:", err)
		return nil, err
	}

	fmt.Println("DB connection opened")
	return conn, nil
}
