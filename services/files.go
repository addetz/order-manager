package jobs

import (
	"encoding/json"
	"log"
	"os"
)

func openJobsFile() []*Job {
	return openFile("data/jobs.json")
}

// func openCustomersFile() []string {
// 	return openFile("data/customers.txt")
// }

func writeJobsFile(rows []*Job) {
	updateFile("data/jobs.json", rows)
}

func openFile(filepath string) []*Job {
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

func updateFile(filepath string, rows []*Job) {
	bytes, err := json.Marshal(rows)
	if err != nil {
		log.Fatal("Error marshal rows:", err)

	}
	os.WriteFile(filepath, bytes, os.ModePerm)
}
