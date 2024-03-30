package jobs

import (
	"encoding/json"
	"log"
	"os"
)

func openJobsFile() []*Job {
	filepath := "data/jobs.json"
	file, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	datas := make([]*Job, 0)
	if err := json.Unmarshal(file, &datas); err != nil {
		log.Fatal("Error unmarshalling file:", err)
	}

	return datas
}

func openCustomersFile() []*Customer {
	filepath := "data/customers.json"
	file, err := os.ReadFile(filepath)
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
	updateFile("data/jobs.json", rows)
}

func updateFile(filepath string, rows []*Job) {
	bytes, err := json.Marshal(rows)
	if err != nil {
		log.Fatal("Error marshal rows:", err)

	}
	os.WriteFile(filepath, bytes, os.ModePerm)
}
