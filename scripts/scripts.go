package main

import (
	_ "embed"

	jobs "github.com/addetz/order-manager/services"
	"honnef.co/go/js/dom"
)

func main() {
	js := jobs.NewJobService()
	document := dom.GetWindow().Document()
	newBody := document.CreateElement("tbody")
	ts := newBody.(*dom.HTMLTableSectionElement)
	for _, e := range js.ListJobs() {
		populateJob(ts, e)
	}
	oldBody := document.GetElementByID("jobsTable").GetElementsByTagName("tbody")[0]
	document.GetElementByID("jobsTable").ReplaceChild(newBody, oldBody)
}

func populateJob(tableSection *dom.HTMLTableSectionElement, job *jobs.Job) {
	row := tableSection.InsertRow(0)
	row.InsertCell(0).SetTextContent(job.ID)
	row.InsertCell(1).SetTextContent(job.OrderDate.Format(jobs.JobsDateFormat))
	row.InsertCell(2).SetTextContent(job.DeadlineDate.Format(jobs.JobsDateFormat))
	row.InsertCell(3).SetTextContent(job.Status)
	row.InsertCell(4).SetTextContent(job.Customer)
	row.InsertCell(5).SetTextContent(job.Description)
}
