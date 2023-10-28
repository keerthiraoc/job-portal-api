package main

import (
	"context"
	"fmt"
	"job-portal-api/internal/auth"
	"job-portal-api/internal/database"
	"job-portal-api/internal/handlers"
	"job-portal-api/internal/services"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func main() {
	err := startApp()
	if err != nil {
		log.Panic().Err(err).Send()
	}
	log.Info().Msg("app end")
}

func startApp() error {
	log.Info().Msg("main : Started : Initializing authentication support")
	privatePEM, err := os.ReadFile("private.pem")
	if err != nil {
		return fmt.Errorf("reading auth private key %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return fmt.Errorf("parsing auth private key %w", err)
	}

	publicPEM, err := os.ReadFile("pubkey.pem")
	if err != nil {
		return fmt.Errorf("reading auth private key %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return fmt.Errorf("parsing auth private key %w", err)
	}

	a, err := auth.NewAuth(publicKey, privateKey)
	if err != nil {
		return fmt.Errorf("constructing auth %w", err)
	}

	log.Info().Msg("main : Started : Initializing db support")
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("connecting to db %w", err)
	}

	pg, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = pg.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("database is not connected: %w", err)
	}

	ms, err := services.NewService(db)
	if err != nil {
		return err
	}

	// err = ms.AutoMigrate()
	// if err != nil {
	// 	return err
	// }

	api := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      handlers.API(ms, a),
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Info().Str("port", api.Addr).Msg("main: API listening")
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error %w", err)
	case sig := <-shutdown:
		log.Info().Msgf("main: Start shutdown %s", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := api.Shutdown(ctx)
		if err != nil {
			err = api.Close()
			return fmt.Errorf("could not stop server gracefully %w", err)
		}
		return nil
	}
}
