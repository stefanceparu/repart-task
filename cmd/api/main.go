package main

import (
	"fmt"
	"log"
	"net/http"
	"reparttask/config"
	"reparttask/internal/order"
	"reparttask/internal/pack"
	"reparttask/service/firstfit"
	"reparttask/storage/memory"
)

func main() {
	router := http.NewServeMux()
	db := memory.NewMemDB()
	calc := firstfit.Calc{}

	packHandler := pack.NewHandler(db)
	packHandler.RegisterRoutes(router)

	orderHandler := order.NewHandler(db, calc)
	orderHandler.RegisterRoutes(router)

	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening on port:", cfg.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
	if err != nil {
		log.Fatal(err)
	}
}
