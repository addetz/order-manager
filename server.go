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
	cs := jobs.NewCustomerService()
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
	// Set up customer view
	e.Static("/customerView", "layout-customers")

	// Set up scripts
	e.File("/scripts.js", "scripts/scripts.js")
	e.File("/scripts.js.map", "scripts/scripts.js.map")
	e.File("/customerView/scripts-customers.js", "scripts-customers/scripts-customers.js")
	e.File("/customerView/scripts-customers.js.map", "scripts-customers/scripts-customers.js.map")

	// List operations
	e.GET("/jobs", func(c echo.Context) error {
		return c.JSON(http.StatusOK, js.ListJobs())
	})

	e.GET("/customers", func(c echo.Context) error {
		return c.JSON(http.StatusOK, cs.ListCustomers())
	})

	// Create operations 
	e.POST("/jobs", func(c echo.Context) error {
		job := &jobs.Job{}
		json.NewDecoder(c.Request().Body).Decode(job)
		js.AddJob(job)
		return c.JSON(http.StatusCreated, nil)
	})

	e.POST("/customers", func(c echo.Context) error {
		cust := &jobs.Customer{}
		json.NewDecoder(c.Request().Body).Decode(cust)
		cs.AddCustomer(cust)
		return c.JSON(http.StatusCreated, nil)
	})

	//Update operations
	e.POST("/jobs/:id", func(c echo.Context) error {
		id := c.Param("id")
		job := &jobs.Job{}
		json.NewDecoder(c.Request().Body).Decode(job)
		js.UpdateJob(id, job)
		return c.JSON(http.StatusOK, nil)
	})

	e.POST("/customers/:id", func(c echo.Context) error {
		id := c.Param("id")
		cust := &jobs.Customer{}
		json.NewDecoder(c.Request().Body).Decode(cust)
		cs.UpdateCustomer(id, cust)
		return c.JSON(http.StatusOK, nil)
	})

	// Delete operations
	e.DELETE("/customers/:id", func(c echo.Context) error {
		id := c.Param("id")
		cs.DeleteCustomer(id)
		return c.JSON(http.StatusOK, nil)
	})

	e.DELETE("/jobs/:id", func(c echo.Context) error {
		id := c.Param("id")
		js.DeleteJob(id)
		return c.JSON(http.StatusOK, nil)
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
