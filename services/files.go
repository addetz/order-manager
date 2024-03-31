package jobs

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

const JOBS_MANAGER_CUSTOMERS_FILE = "jobsManager-customers.json"
const JOBS_MANAGER_JOBS_FILE = "jobsManager-jobs.json"

func openJobsFile() []*Job {
	if _, err := os.Stat(JOBS_MANAGER_JOBS_FILE); errors.Is(err, os.ErrNotExist) {
		createEmptyFile(JOBS_MANAGER_JOBS_FILE)
		writeJobsFile([]*Job{})
		return []*Job{}
	}
	file, err := os.ReadFile(JOBS_MANAGER_JOBS_FILE)
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

func openCustomersFile() []*Customer {
	if _, err := os.Stat(JOBS_MANAGER_CUSTOMERS_FILE); errors.Is(err, os.ErrNotExist) {
		createEmptyFile(JOBS_MANAGER_CUSTOMERS_FILE)
		writeCustomersFile([]*Customer{})
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

func writeJobsFile(rows []*Job) {
	bytes, err := json.Marshal(rows)
	if err != nil {
		log.Fatal("Error marshal rows:", err)

	}
	os.WriteFile(JOBS_MANAGER_JOBS_FILE, bytes, os.ModePerm)
}

func writeCustomersFile(rows []*Customer) {
	bytes, err := json.Marshal(rows)
	if err != nil {
		log.Fatal("Error marshal rows:", err)

	}
	os.WriteFile(JOBS_MANAGER_CUSTOMERS_FILE, bytes, os.ModePerm)
}
