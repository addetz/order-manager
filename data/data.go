package data

import (
	"bufio"
	"log"
	"os"
)

func OpenJobsFile() []string {
	return openFile("./data/jobs.txt")
}

func OpenCustomersFile() []string {
	return openFile("./data/customers.txt")
}

func WriteJobsFile(rows []string) {
	updateFile("./data/jobs.txt", rows)
}

func openFile(filepath string) []string {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func updateFile(filepath string, rows []string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	for _, r := range rows {
		file.WriteString(r)
	}
}
