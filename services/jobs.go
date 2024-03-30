package jobs

import (
	_ "embed"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

const JobsDateFormat string = "Mon 02 Jan 2006"

//go:embed jobs
var jobsInput string

var JobStatusList []string = []string{
	"New ⭐️",
	"Completed & Shipped✅",
	"Invoiced 🧾",
}

type Job struct {
	ID           string
	OrderDate    time.Time
	DeadlineDate time.Time
	Status       string
	Customer     string
	Description  string
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

func (js *JobService) AddJob(orderDate string, deadline string, customer string,
	description string) {
	id := fmt.Sprintf("#%d", len(js.jobs)+1)
	order := &Job{
		ID:           id,
		OrderDate:    *getFormattedDate(orderDate),
		DeadlineDate: *getFormattedDate(deadline),
		Status:       JobStatusList[0],
		Customer:     customer,
		Description:  description,
	}
	js.jobs[order.ID] = order
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
	rows := strings.Split(jobsInput, "\n")
	for _, row := range rows {
		cells := strings.Split(row, ",")
		job := &Job{
			ID:           cells[0],
			OrderDate:    *getFormattedDate(cells[1]),
			DeadlineDate: *getFormattedDate(cells[2]),
			Status:       cells[3],
			Customer:     cells[4],
			Description:  cells[5],
		}
		js.jobs[job.ID] = job
	}
}

func getFormattedDate(s string) *time.Time {
	t, err := time.Parse(JobsDateFormat, s)
	if err != nil {
		log.Fatal(err)
	}
	return &t
}
