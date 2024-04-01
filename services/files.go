package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

const JOBS_MANAGER_CUSTOMERS_FILE = "jobsManager-customers.json"
const JOBS_MANAGER_JOBS_FILE = "jobsManager-jobs.json"

func openJobsFile(filepath string) []*Job {
	fullPath := fmt.Sprintf("%s/%s", filepath, JOBS_MANAGER_JOBS_FILE)
	if _, err := os.Stat(fullPath); errors.Is(err, os.ErrNotExist) {
		createEmptyFile(fullPath)
		writeJobsFile(filepath, []*Job{})
		return []*Job{}
	}
	file, err := os.ReadFile(fullPath)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	datas := make([]*Job, 0)
	if err := json.Unmarshal(file, &datas); err != nil {
		log.Fatal("Error unmarshalling file:", err)
	}

	return datas
}

func createEmptyFile(title string) {
	file, err := os.Create(title) //create a new file
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
}

func openCustomersFile(filepath string) []*Customer {
	fullPath := fmt.Sprintf("%s/%s", filepath, JOBS_MANAGER_CUSTOMERS_FILE)
	if _, err := os.Stat(JOBS_MANAGER_CUSTOMERS_FILE); errors.Is(err, os.ErrNotExist) {
		createEmptyFile(fullPath)
		writeCustomersFile(filepath, []*Customer{})
		return []*Customer{}
	}
	file, err := os.ReadFile(JOBS_MANAGER_CUSTOMERS_FILE)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	datas := make([]*Customer, 0)
	if err := json.Unmarshal(file, &datas); err != nil {
		log.Fatal("Error unmarshalling file:", err)
	}

	return datas
}

func writeJobsFile(filepath string, rows []*Job) {
	fullPath := fmt.Sprintf("%s/%s", filepath, JOBS_MANAGER_JOBS_FILE)
	bytes, err := json.Marshal(rows)
	if err != nil {
		log.Fatal("Error marshal rows:", err)

	}
	os.WriteFile(fullPath, bytes, os.ModePerm)
}

func writeCustomersFile(filepath string, rows []*Customer) {
	fullPath := fmt.Sprintf("%s/%s", filepath, JOBS_MANAGER_CUSTOMERS_FILE)
	bytes, err := json.Marshal(rows)
	if err != nil {
		log.Fatal("Error marshal rows:", err)

	}
	os.WriteFile(fullPath, bytes, os.ModePerm)
}
