package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jobs "github.com/addetz/order-manager/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed frontend/layout/index.html
var rootIndex []byte

//go:embed frontend/layoutCustomers/index.html
var customerIndex []byte

//go:embed frontend/scripts/scripts.js
var scripts []byte

//go:embed frontend/scripts/scripts.js.map
var scriptsMap []byte

//go:embed frontend/scriptsCustomers/scriptsCustomers.js
var scriptsCustomers []byte

//go:embed frontend/scriptsCustomers/scriptsCustomers.js.map
var scriptsCustomersMap []byte

//go:embed frontend/layout/custom.css
var customCSS []byte

//go:embed frontend/layout/favicon-melon.ico
var favicon []byte

const (
	TIMEOUT = 3 * time.Second
)

func main() {
	filePath := flag.String("filepath", ".", "executable path")
	cs := jobs.NewCustomerService(*filePath)
	js := jobs.NewJobService(*filePath)

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
	e.GET("/", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "text/html; charset=utf-8", rootIndex)
	})

	e.GET("/favicon-melon.ico", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "image/x-icon", favicon)
	})

	e.GET("/customerView", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "text/html; charset=utf-8", customerIndex)
	})

	e.GET("/scripts.js", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/javascript", scripts)
	})

	e.GET("/custom.css", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "text/css; charset=utf-8", customCSS)
	})

	e.GET("/scripts.js.map", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/javascript", scriptsMap)
	})

	e.GET("/scriptsCustomers.js", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/javascript", scriptsCustomers)
	})

	e.GET("/scriptsCustomers.js.map", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/javascript", scriptsCustomersMap)
	})

	// List operations
	e.GET("/jobs", func(c echo.Context) error {
		customerID := c.QueryParam("customerID")
		log.Println(customerID)
		if customerID == "unknown" {
			return c.JSON(http.StatusOK, js.FilterJobs(""))
		}
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
		job := jobs.NewJob("", "", "", "", "")
		json.NewDecoder(c.Request().Body).Decode(job)
		log.Printf("\n\n%v\n\n", job)
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
