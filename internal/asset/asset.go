package asset

import (
	"github.com/google/uuid"
)

type Asset struct {
	ID         string
	Name       string
	Ticker     string
	IsTradable bool
}

func NewAsset(name string, ticker string, IsTradable bool) *Asset {
	return &Asset{
		ID:         uuid.New().String(),
		Name:       name,
		Ticker:     ticker,
		IsTradable: IsTradable,
	}
}

func DisableAsset(ticker string) {
	// disabling trading logic.
}
