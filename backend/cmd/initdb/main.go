package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := "host=localhost port=5432 user=postgres password=postgres sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer pool.Close()

	var exists bool
	err = pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'cashier_config')").Scan(&exists)
	if err != nil {
		log.Fatalf("failed to check db: %v", err)
	}
	if exists {
		fmt.Println("database cashier_config already exists")
		return
	}

	_, err = pool.Exec(context.Background(), "CREATE DATABASE cashier_config")
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	fmt.Println("database cashier_config created")
}
