package main

import (
	"net/http"

	"little-farm-webserver/internal/infrastructure/controller"
	"little-farm-webserver/internal/infrastructure/persistence"
)

func main() {
	mux := http.NewServeMux()

	staticDir := http.Dir("./static")
	staticServer := http.FileServer(staticDir)
	mux.Handle("GET /", http.StripPrefix("/", staticServer))

	farmDAO := persistence.NewInMemoryFarmDAO()
	farmController := controller.NewFarmController(farmDAO)
	mux.HandleFunc("POST /api/farms", farmController.CreateFarm)
	mux.HandleFunc("GET /api/farms/{farm_id}", farmController.FindFarm)
	mux.HandleFunc("POST /api/farms/{farm_id}/advance-one-second", farmController.AdvanceOneSecond)
	mux.HandleFunc("POST /api/farms/{farm_id}/corn/{corn_id}/harvest", farmController.HarvestCorn)
	mux.HandleFunc("POST /api/farms/{farm_id}/grass/{grass_id}/harvest", farmController.HarvestGrass)
	mux.HandleFunc("POST /api/farms/{farm_id}/dirt/{dirt_id}/plant-corn", farmController.PlantCorn)
	mux.HandleFunc("POST /api/farms/{farm_id}/dirt/{dirt_id}/plant-wheat", farmController.PlantWheat)

	(&http.Server{Addr: ":8080", Handler: mux}).ListenAndServe()
}
