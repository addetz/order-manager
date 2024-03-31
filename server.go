package main

import (
	"embed"
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

//go:embed frontend
var frontend embed.FS

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
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: http.FS(frontend),
	}))

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
	e.GET("/", func(c echo.Context) error {
		return c.File("frontend/layout/index.html")
	})

	e.GET("/customerView", func(c echo.Context) error {
		return c.File("frontend/layoutCustomers/index.html")
	})

	e.GET("/scripts.js", func(c echo.Context) error {
		return c.File("frontend/scripts/scripts.js")
	})

	e.GET("/custom.css", func(c echo.Context) error {
		return c.File("frontend/layout/custom.css")
	})

	e.GET("/scripts.js.map", func(c echo.Context) error {
		return c.File("frontend/scripts/scripts.js.map")
	})

	e.GET("/customerView/scriptsCustomers.js", func(c echo.Context) error {
		return c.File("frontend/scriptsCustomers/scriptsCustomers.js")
	})

	e.GET("/customerView/scriptsCustomers.js.map", func(c echo.Context) error {
		return c.File("frontend/scriptsCustomers/scriptsCustomers.js.map")
	})

	// List operations
	e.GET("/jobs", func(c echo.Context) error {
		customerID := c.QueryParam("customerID")
		if customerID != "" {
			return c.JSON(http.StatusOK, js.FilterJobs(customerID))
		}
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

	//Search customer
	e.GET("/customers/search", func(c echo.Context) error {
		customerName := c.QueryParam("name")
		customer, err := cs.SearchCustomer(customerName)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusOK, customer)
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
