package main

import (
	"DriveApi/internal/app"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializting configs: %s", err.Error())
	}
	srv := app.NewServer()
	srv.Logger.Fatal(srv.Start())
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
