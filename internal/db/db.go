package db

import (
	"database/sql"
	"fmt"
	"log"

	go_ora "github.com/sijms/go-ora/v2"
	"saranasistemsolusindo.com/gusen-admin/internal/config"
)

// DBConfig holds the database configuration parameters
type DBConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	ServiceName string
}

// InitializeDB establishes a connection to the Oracle database
func InitializeDB() (*sql.DB, error) {
	port := config.OracleInfo.Port
	host := config.OracleInfo.Host
	serviceName := config.OracleInfo.ServiceName
	user := config.OracleInfo.User
	pass := config.OracleInfo.Pass

	if host == "" || user == "" || pass == "" || serviceName == "" || port == 0 {
		return nil, fmt.Errorf("database configuration is incomplete")
	}

	connStr := go_ora.BuildUrl(host, port, serviceName, user, pass, nil)

	db, err := sql.Open("oracle", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Successfully connected to the database")

	// Configure connection pool settings if needed
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(5)
	// db.SetConnMaxLifetime(time.Hour)

	return db, nil
}
