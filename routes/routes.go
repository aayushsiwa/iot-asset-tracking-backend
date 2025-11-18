package routes

import (
	"crud/handlers"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRoutes() Routes {
	return Routes{
		// Health Check
		{
			Name:    "HealthCheck",
			Method:  http.MethodGet,
			Pattern: "/health",
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("ok"))
			},
		},
		// Locations
		{
			Name:        "CreateLocation",
			Method:      http.MethodPost,
			Pattern:     "/locations",
			HandlerFunc: handlers.CreateLocation,
		},
		{
			Name:        "GetLocation",
			Method:      http.MethodGet,
			Pattern:     "/locations",
			HandlerFunc: handlers.GetLocation,
		},
		{
			Name:        "UpdateLocation",
			Method:      http.MethodPatch,
			Pattern:     "/locations/{id}",
			HandlerFunc: handlers.UpdateLocation,
		},
		{
			Name:        "DeleteLocation",
			Method:      http.MethodDelete,
			Pattern:     "/locations/{id}",
			HandlerFunc: handlers.DeleteLocation,
		},
		// Assets
		{
			Name:        "CreateAsset",
			Method:      http.MethodPost,
			Pattern:     "/locations/{locationID}/assets",
			HandlerFunc: handlers.CreateAsset,
		},
		{
			Name:        "GetAssetsByLocation",
			Method:      http.MethodGet,
			Pattern:     "/locations/{locationID}/assets",
			HandlerFunc: handlers.GetAssetsByLocation,
		},
		{
			Name:        "GetAssets",
			Method:      http.MethodGet,
			Pattern:     "/assets",
			HandlerFunc: handlers.GetAssets,
		},
		{
			Name:        "UpdateAssets",
			Method:      http.MethodPatch,
			Pattern:     "/locations/{locationID}/assets/{assetID}",
			HandlerFunc: handlers.UpdateAsset,
		},
		{
			Name:        "DeleteAsset",
			Method:      http.MethodDelete,
			Pattern:     "/locations/{locationID}/assets/{assetID}",
			HandlerFunc: handlers.DeleteAsset,
		},
	}
}

func AttachRoutes(router *Router, routes Routes) {
	for _, route := range routes {
		router.Handle(route.Method, route.Pattern, route.HandlerFunc)
	}
}
