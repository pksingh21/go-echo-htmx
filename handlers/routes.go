package handlers

import "github.com/labstack/echo/v4"

func SetupRoutes(app *echo.Echo, h *UserHandler) {
	group := app.Group("/user")
	group.GET("", h.LearningHandler)
	group.GET("/details/:id", h.HandlerShowUserById)
	group.GET("/info", h.HandlerShowUsers)
}
