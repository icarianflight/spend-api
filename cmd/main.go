package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	dbAccounts "spend-api/internal/app/adapters/db/accounts"
	restAccounts "spend-api/internal/app/adapters/rest/accounts"
	"spend-api/internal/config"
	domainAccounts "spend-api/internal/domain/accounts"
	"spend-api/internal/infra/db"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {

	cfg := config.LoadConfig()

	executor, err := db.NewMariaDbExecutor(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func(executor *db.MariaDbExecutor) {
		_ = executor.Close()
	}(executor)

	accountDbAdapter := dbAccounts.NewForSavingAccountUsingDB(executor)

	accountService := domainAccounts.NewAccountService(accountDbAdapter)

	accountAPIHandler := restAccounts.NewForCreatingAccountUsingRestAPI(accountService)

	http.Handle("/accounts", accountAPIHandler)
	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
