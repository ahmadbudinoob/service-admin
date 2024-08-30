package router

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"saranasistemsolusindo.com/gusen-admin/internal/handlers"
	"saranasistemsolusindo.com/gusen-admin/internal/utils/jwt"
)

// InitRouter initializes the Echo router with the defined routes
func InitRouter(db *sql.DB) (*echo.Echo, error) {
	e := echo.New()

	// Initialize UserHandler
	userHandler, err := handlers.NewUserHandler(db)
	if err != nil {
		return nil, err
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Public routes
	e.POST("/login", userHandler.LoginAdmin)

	// Group routes that require JWT authentication
	authGroup := e.Group("", echo.WrapMiddleware(jwt.JWTMiddleware))
	authGroup.GET("/paginated", userHandler.GetUserPaginated)
	authGroup.POST("/create", userHandler.CreateUser)
	authGroup.PUT("/update/:id", userHandler.UpdateUser)
	authGroup.GET("/user/:login_id", userHandler.GetUserByLoginId)
	authGroup.GET("/log_history", userHandler.GetLogHistory)
	authGroup.GET("/client/login/:login_id", userHandler.GetClientByLoginID)
	authGroup.GET("/client/", userHandler.GetClientByClientID)
	authGroup.PUT("/client/update", userHandler.UpdateClientByUserLogin)
	authGroup.GET("/client/not_in_user", userHandler.GetAvailableClients)
	authGroup.PUT("/deactive-user/:login_id", userHandler.DeactiveUser)
	// reset pin and password
	authGroup.PUT("/reset-pin", userHandler.ResetPin)
	authGroup.PUT("/reset-password", userHandler.ResetPassword)
	return e, nil
}
