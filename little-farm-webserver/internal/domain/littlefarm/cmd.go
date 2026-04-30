package littlefarm

import (
	"time"

	"little-farm-webserver/pkg/domain"
)

type CreateFarmCmd struct {
	FarmID    domain.ID
	MapWidth  int
	MapHeight int
	Timestamp time.Time
}

type FarmCreatedEvt struct {
	FarmID          domain.ID
	MapWidth        int
	MapHeight       int
	Resources       Resources
	ActiveElements  []ActiveElement
	PassiveElements []ActiveElement
	Timestamp       time.Time
}

type AdvanceOneSecondCmd struct {
	FarmID    domain.ID
	Timestamp time.Time
}

type OneSecondAdvancedEvt struct {
	FarmID    domain.ID
	Timestamp time.Time
}

type HarvestCornCmd struct {
	FarmID    domain.ID
	CornID    domain.ID
	Timestamp time.Time
}

type CornHarvestedEvt struct {
	FarmID    domain.ID
	CornID    domain.ID
	Produced  int
	Timestamp time.Time
}

type HarvestGrassCmd struct {
	FarmID    domain.ID
	GrassID   domain.ID
	Timestamp time.Time
}

type GrassHarvestedEvt struct {
	FarmID    domain.ID
	GrassID   domain.ID
	Produced  int
	Timestamp time.Time
}
