package model

import (
	"time"

	"github.com/google/uuid"
)

type Status string

// Statuses is a map of statuses
var Statuses = struct {
	Offline Status
	Online  Status
	values  map[Status]bool
}{
	Offline: "offline",
	Online:  "online",
}

type Asset struct {
	ID               *uuid.UUID `json:"ID"`
	Name             string     `json:"name"`
	Status           Status     `json:"status"`
	Location         string     `json:"location"`
	LastUpdatedAtUTC time.Time  `json:"lastUpdatedAtUTC"`
	CreatedAtUTC     time.Time  `json:"createdAtUTC"`
}

type AssetPatch struct {
	Name   *string
	Status *Status
}

type CreateAssetRequest struct {
	ID         *uuid.UUID `json:"ID"`
	Name       string     `json:"name"`
	Status     Status     `json:"status"`
	LocationID uuid.UUID  `json:"locationID"`
}
