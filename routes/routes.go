package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"myapp/controllers"
	middlewareCustom "myapp/middleware"
)

func SetupRoutes(e *echo.Echo) {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize controllers
	authController := &controllers.AuthController{}
	attendanceController := &controllers.AttendanceController{}

	// Public routes
	api := e.Group("/api/v1")
	
	// Health check
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "OK",
			"message": "Attendance System API is running",
		})
	})

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.POST("/login", authController.Login)
	auth.POST("/register", authController.Register)
	auth.POST("/refresh", authController.RefreshToken)

	// Protected routes (require JWT)
	protected := api.Group("")
	protected.Use(middlewareCustom.JWTMiddleware())

	// User profile routes
	protected.GET("/profile", authController.GetProfile)

	// Attendance routes (protected)
	attendanceRoutes := protected.Group("/attendance")
	attendanceRoutes.POST("/record", attendanceController.RecordAttendance)
	attendanceRoutes.GET("/today", attendanceController.GetTodayAttendance)
	attendanceRoutes.GET("/history/:student_id", attendanceController.GetAttendanceHistory)

	// Admin routes (require admin role)
	admin := protected.Group("/admin")
	admin.Use(middlewareCustom.AdminMiddleware())
	admin.POST("/nfc/register", attendanceController.RegisterNFCCard)

	// Super admin routes (require super admin role)
	superAdmin := protected.Group("/super-admin")
	superAdmin.Use(middlewareCustom.SuperAdminMiddleware())
	// Add super admin specific routes here
	superAdmin.GET("/users", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "Super admin users endpoint",
		})
	})
}