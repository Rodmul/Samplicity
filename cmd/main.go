package main

import (
	"DriveApi/internal/app"
	"DriveApi/internal/repository"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializting configs: %s", err.Error())
	}

	_, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: "1234",
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	log.Println("Connect to database postgres")

	srv := app.NewServer()
	srv.Logger.Fatal(srv.Start())
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
