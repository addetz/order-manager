package jobs

import (
	_ "embed"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	b64 "encoding/base64"

	"github.com/addetz/order-manager/data"
)

const JobsDateFormat string = "2006-01-02"

var JobStatusList []string = []string{
	"New ⭐️",
	"Completed & Shipped✅",
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
	rows := data.OpenJobsFile()
	for _, row := range rows {
		cells := strings.Split(row, ",")
		decodedDescription, _ := b64.StdEncoding.DecodeString(cells[5])
		job := &Job{
			ID:           cells[0],
			OrderDate:    *getFormattedDate(cells[1]),
			DeadlineDate: *getFormattedDate(cells[2]),
			Status:       cells[3],
			Customer:     cells[4],
			Description:  string(decodedDescription),
		}
		js.jobs[job.ID] = job
	}
}

func (js *JobService) exportJobs() {
	jobsList := js.ListJobs()
	rows := make([]string, len(jobsList))
	for i, j := range jobsList {
		orderDate := j.OrderDate.Format(JobsDateFormat)
		deadlineDate := j.DeadlineDate.Format(JobsDateFormat)
		description := b64.StdEncoding.EncodeToString([]byte(j.Description))
		rows[i] = fmt.Sprintf("%s,%s,%s,%s,%s,%s", j.ID, orderDate, deadlineDate,
			j.Status, j.Customer, description)
	}
}

func getFormattedDate(s string) *time.Time {
	t, err := time.Parse(JobsDateFormat, s)
	if err != nil {
		log.Fatal(err)
	}
	return &t
}
