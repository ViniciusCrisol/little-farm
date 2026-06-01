package littlefarm

import (
	"fmt"
	"time"

	"little-farm-webserver/pkg/apperr"
	"little-farm-webserver/pkg/domain"
)

var (
	ErrInvalidFarmID    = fmt.Errorf("%w: invalid farm id", apperr.ErrValidation)
	ErrFarmNotFound     = fmt.Errorf("%w: farm not found", apperr.ErrValidation)
	ErrInvalidMapWidth  = fmt.Errorf("%w: invalid map width", apperr.ErrValidation)
	ErrInvalidMapHeight = fmt.Errorf("%w: invalid map height", apperr.ErrValidation)

	ErrInvalidCornID = fmt.Errorf("%w: invalid corn id", apperr.ErrValidation)
	ErrCornNotFound  = fmt.Errorf("%w: corn not found", apperr.ErrValidation)
	ErrCornNotReady  = fmt.Errorf("%w: corn not ready", apperr.ErrValidation)

	ErrInvalidGrassID = fmt.Errorf("%w: invalid grass id", apperr.ErrValidation)
	ErrGrassNotFound  = fmt.Errorf("%w: grass not found", apperr.ErrValidation)
	ErrGrassNotReady  = fmt.Errorf("%w: grass not ready", apperr.ErrValidation)

	ErrInvalidWheatID = fmt.Errorf("%w: invalid wheat id", apperr.ErrValidation)
	ErrWheatNotFound  = fmt.Errorf("%w: wheat not found", apperr.ErrValidation)
	ErrWheatNotReady  = fmt.Errorf("%w: wheat not ready", apperr.ErrValidation)
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

type AdvanceOneSecondCmd struct {
	FarmID    domain.ID
	Timestamp time.Time
}

type OneSecondAdvancedEvt struct {
	FarmID    domain.ID
	Timestamp time.Time
}
