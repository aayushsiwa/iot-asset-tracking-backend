package routes

import (
	"crud/middleware"
)

func SetupRouter() *Router {
	router := NewRouter()

	// middlewares
	router.Use(middleware.RecoverMiddleware)
	router.Use(middleware.LoggingMiddleware)

	// group
	api := router.Group("/api/v1")

	// routes
	apiRoutes := NewRoutes()
	AttachRoutes(api, apiRoutes)

	return router
}
