package main

import (
	"encoding/json"
	"fmt"
	"time"

	"little-farm-webserver/internal/domain/littlefarm"
	"little-farm-webserver/pkg/domain"
)

func main() {
	farm, err := littlefarm.NewFarm(littlefarm.CreateFarmCmd{
		FarmID:    domain.GenerateID(),
		MapWidth:  16,
		MapHeight: 16,
		Timestamp: time.Now(),
	})
	if err != nil {
		panic(err)
	}

	j, err := json.Marshal(farm.UncommittedEvents())
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
}
