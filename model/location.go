package model

import "github.com/google/uuid"

type Location struct {
	ID               *uuid.UUID `json:"ID"`
	Name             string     `json:"name"`
	Code             string     `json:"code"`
	CreatedAtUTC     string     `json:"createdAtUTC"`
	LastUpdatedAtUTC string     `json:"lastUpdatedAtUTC"`
}

type UpdateLocationRequest struct {
	Name *string `json:"name" validate:"omitempty,min=5,max=50"`
	Code *string `json:"code" validate:"omitempty,uppercase,len=4"`
}

type LocationPatch struct {
	ID   uuid.UUID
	Name *string
	Code *string
}
