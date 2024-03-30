package jobs

import (
	_ "embed"
	"fmt"
	"log"
	"sort"
	"time"
)

const JobsDateFormat string = "2006-01-02"

var JobStatusList []string = []string{
	"New ⭐️",
	"Completed & Shipped ✅",
	"Invoiced 🧾",
}

// TODO remove hardcoding
var CustomerList []string = []string{
	"Adelina",
	"Stuzzlini",
}

type Job struct {
	ID           string    `json:"id"`
	OrderDate    time.Time `json:"order_date"`
	DeadlineDate time.Time `json:"deadline_date"`
	Status       string    `json:"status"`
	Customer     string    `json:"customer"`
	Description  string    `json:"description"`
}

type JobService struct {
	jobs map[string]*Job
}

func NewJobService() *JobService {
	js := &JobService{
		jobs: make(map[string]*Job, 0),
	}
	js.importJobs()
	return js
}

func (js *JobService) AddJob(j *Job) {
	id := fmt.Sprintf("#%d", len(js.jobs)+1)
	j.ID = id
	js.jobs[id] = j
	js.exportJobs()
}

func (js *JobService) ListJobs() []*Job {
	jobsList := make([]*Job, 0)
	for _, o := range js.jobs {
		jobsList = append(jobsList, o)
	}
	sort.Slice(jobsList, func(i, j int) bool {
		return jobsList[i].ID > jobsList[j].ID
	})
	return jobsList
}

func (js *JobService) importJobs() {
	rows := openJobsFile()
	for _, row := range rows {
		js.jobs[row.ID] = row
	}
}

func (js *JobService) exportJobs() {
	writeJobsFile(js.ListJobs())
}

func getFormattedDate(s string) *time.Time {
	t, err := time.Parse(JobsDateFormat, s)
	if err != nil {
		log.Fatal(err)
	}
	return &t
}
