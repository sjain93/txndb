package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sjain93/userservice/api/txns"
	"github.com/sjain93/userservice/api/user"
	"github.com/sjain93/userservice/config"
	"github.com/sjain93/userservice/migrations"
	"github.com/sjain93/userservice/routes"
)

func main() {
	var (
		userRepository user.UserRepoManager
		txnRepository  txns.TxnRepoManager
		err            error
	)

	err = config.LoadEnvVars()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	config.ConnectDatabase()
	migrations.AutoMigrate(config.DB)

	userRepository, err = user.NewUserRepository(config.DB)
	if err != nil {
		log.Fatalf("Error initializing user repo: %v", err.Error())
	}

	userService := user.NewUserService(userRepository)

	txnRepository, err = txns.NewTxnRepository(config.DB)
	if err != nil {
		log.Fatalf("Error initializing txn repo: %v", err.Error())
	}

	txnService := txns.NewTxnService(txnRepository)
	if err := txnService.SeedLocalTransactions(); err != nil {
		log.Fatalf("Error seeding txn repo: %v", err.Error())
	}

	e := echo.New()
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.RequestID(),
	)

	routes.SetupAPIRoutes(e, userService)

	/*
		The code below implements a graceful shutdown by starting the server
		via  a goroutine that blocks until a kill command is posted
	*/

	shutdownCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-shutdownCtx.Done() // block here until ctrl+c

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

/*
	An optional route can be added to trigger a graceful shutdown over HTTP:

	e.POST("/quit", func(c echo.Context) error {
		cancel()
		return c.String(http.StatusOK, "OK")
	})
*/
