package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jobs "github.com/addetz/order-manager/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	TIMEOUT = 3 * time.Second
)

func main() {
	js := jobs.NewJobService()

	// Read port if one is set
	port := readPort()

	// Initialise echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Configure server
	s := http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           e,
		ReadTimeout:       TIMEOUT,
		ReadHeaderTimeout: TIMEOUT,
		WriteTimeout:      TIMEOUT,
		IdleTimeout:       TIMEOUT,
	}

	// Set up the root file
	e.Static("/", "layout")
	// Set up scripts
	e.File("/scripts.js", "scripts/scripts.js")
	e.File("/scripts.js.map", "scripts/scripts.js.map")
	e.File("/favicon.ico", "layout/favicon.ico")

	e.GET("/jobs", func(c echo.Context) error {
		jobs := js.ListJobs()
		return c.JSON(http.StatusOK, jobs)
	})

	e.POST("/jobs", func(c echo.Context) error {
		job := &jobs.Job{}
		json.NewDecoder(c.Request().Body).Decode(job)
		js.AddJob(job)
		return c.JSON(http.StatusCreated, nil)
	})

	log.Printf("Listening on localhost:%s...\n", port)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

// readPort reads the SERVER_PORT environment variable if one is set
// or returns a default if none is found
func readPort() string {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		return "8080"
	}
	return port
}
