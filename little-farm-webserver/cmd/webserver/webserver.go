package main

import (
	"net/http"

	"little-farm-webserver/internal/infrastructure/controller"
	"little-farm-webserver/internal/infrastructure/persistence"
)

func main() {
	dao := persistence.NewInMemoryFarmDAO()
	controller := controller.NewFarmController(dao)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /farms", controller.CreateFarm)
	mux.HandleFunc("GET /farms/{farm_id}", controller.FindFarm)
	mux.HandleFunc("POST /farms/{farm_id}/advance-one-second", controller.AdvanceOneSecond)
	mux.HandleFunc("POST /farms/{farm_id}/corn/{corn_id}/harvest", controller.HarvestCorn)
	mux.HandleFunc("POST /farms/{farm_id}/grass/{grass_id}/harvest", controller.HarvestGrass)

	(&http.Server{Addr: ":8080", Handler: mux}).ListenAndServe()
}
