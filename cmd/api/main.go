package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gastro-go/internal/database"
	"gastro-go/internal/handler"
	"gastro-go/internal/repository"
	"gastro-go/internal/usecase"
)

func main() {
	// Initialize database connection
	ctx := context.Background()
	pool, err := database.NewConnection(ctx)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Initialize SQLC queries
	queries := database.New(pool)

	// Initialize repository
	restaurantRepo := repository.NewRestaurantRepository(queries)

	// Initialize use cases
	createRestaurantUC := usecase.NewCreateRestaurantUseCase(restaurantRepo)
	listRestaurantsUC := usecase.NewListRestaurantsUseCase(restaurantRepo)
	getRestaurantBySlugUC := usecase.NewGetRestaurantBySlugUseCase(restaurantRepo)
	openRestaurantUC := usecase.NewOpenRestaurantUseCase(restaurantRepo)
	closeRestaurantUC := usecase.NewCloseRestaurantUseCase(restaurantRepo)
	updateOpeningHoursUC := usecase.NewUpdateOpeningHoursUseCase(restaurantRepo)
	updatePaymentMethodsUC := usecase.NewUpdatePaymentMethodsUseCase(restaurantRepo)

	// Initialize handler
	restaurantHandler := handler.NewRestaurantHandler(
		createRestaurantUC,
		listRestaurantsUC,
		getRestaurantBySlugUC,
		openRestaurantUC,
		closeRestaurantUC,
		updateOpeningHoursUC,
		updatePaymentMethodsUC,
	)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "gastro-go",
		})
	})

	// Restaurant routes
	e.POST("/restaurants", restaurantHandler.CreateRestaurant)
	e.GET("/restaurants", restaurantHandler.ListRestaurants)
	e.GET("/restaurants/:slug", restaurantHandler.GetRestaurantBySlug)
	e.PATCH("/restaurants/:id/open", restaurantHandler.OpenRestaurant)
	e.PATCH("/restaurants/:id/close", restaurantHandler.CloseRestaurant)
	e.PUT("/restaurants/:id/hours", restaurantHandler.UpdateOpeningHours)
	e.PUT("/restaurants/:id/payments", restaurantHandler.UpdatePaymentMethods)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server in a goroutine
	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutting down the server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server gracefully stopped")
}

