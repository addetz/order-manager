package main

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	jobs "github.com/addetz/order-manager/services"
	"honnef.co/go/js/dom"
)

func main() {
	document := dom.GetWindow().Document()
	populateAllJobs(document)

	addRowBtn := document.GetElementByID("addRowBtn")
	submitBtn := document.GetElementByID("submitBtn")
	cancelBtn := document.GetElementByID("cancelBtn")
	addRowBtn.AddEventListener("click", true, func(e dom.Event) {
		showUserInput(document)
	})
	submitBtn.AddEventListener("click", true, func(e dom.Event) {
		submitJob(document)
	})
	cancelBtn.AddEventListener("click", true, func(e dom.Event) {
		hideUserInput(document)
	})

	statusDropdown := document.GetElementByID("statusDropdown").(*dom.HTMLSelectElement)
	populateStatusDropdownOptions(document, statusDropdown, "")
	customerDropdown := document.GetElementByID("customerDropdown").(*dom.HTMLSelectElement)
	populateCustomerDropdownOptions(document, customerDropdown, "")
}

func populateStatusDropdownOptions(document dom.Document,
	statusDropdown *dom.HTMLSelectElement,
	currentValue string) {
	for i, c := range jobs.JobStatusList {
		o := document.CreateElement("option")
		o.SetTextContent(c)
		statusDropdown.AppendChild(o)
		if c == currentValue {
			statusDropdown.SelectedIndex = i
		}
	}
}

func populateCustomerDropdownOptions(document dom.Document,
	customerDropdown *dom.HTMLSelectElement,
	currentValue string) {
	go func(callback func(document dom.Document,
		customerDropdown *dom.HTMLSelectElement,
		customers []*jobs.Customer,
		currentValue string)) {
		resp, err := http.Get("/customers")
		if err != nil {
			log.Fatal(err)
		}
		customers, err := jobs.NewCustomersResponse(resp)
		if err != nil {
			log.Fatal(err)
		}
		callback(document, customerDropdown, customers, currentValue)
	}(populateCustomerDropdownOptionsCallback)
}

func populateCustomerDropdownOptionsCallback(document dom.Document,
	customerDropdown *dom.HTMLSelectElement,
	customers []*jobs.Customer,
	currentValue string) {
	for i, c := range customers {
		o := document.CreateElement("option")
		o.SetTextContent(c.Name)
		customerDropdown.AppendChild(o)
		if c.Name == currentValue {
			customerDropdown.SelectedIndex = i
		}
	}
}

func populateAllJobs(document dom.Document) {
	go func(callback func(document dom.Document, jobs []*jobs.Job)) {
		resp, err := http.Get("/jobs")
		if err != nil {
			log.Fatal(err)
		}
		jobs, err := jobs.NewJobsResponse(resp)
		if err != nil {
			log.Fatal(err)
		}
		callback(document, jobs)
	}(populateJobsCallback)
}

func populateJobsCallback(document dom.Document, jobs []*jobs.Job) {
	newBody := document.CreateElement("tbody")
	ts := newBody.(*dom.HTMLTableSectionElement)
	for _, e := range jobs {
		populateJob(document, ts, e)
	}
	oldBody := document.GetElementByID("jobsTable").GetElementsByTagName("tbody")[0]
	document.GetElementByID("jobsTable").ReplaceChild(newBody, oldBody)
}

func populateJob(document dom.Document, tableSection *dom.HTMLTableSectionElement, job *jobs.Job) {
	row := tableSection.InsertRow(0)
	row.InsertCell(0).SetTextContent(job.ID)
	row.InsertCell(1).SetTextContent(job.OrderDate.Format(jobs.JobsDateFormat))
	row.InsertCell(2).SetTextContent(job.DeadlineDate.Format(jobs.JobsDateFormat))

	// Status
	statusCell := row.InsertCell(3)
	statusCell.SetContentEditable("true")
	statusSelectElement := document.CreateElement("select").(*dom.HTMLSelectElement)
	statusSelectElement.Class().Add("form-control")
	statusSelectElement.SetID(fmt.Sprintf("statusDropdown-%s", job.ID))
	populateStatusDropdownOptions(document, statusSelectElement, job.Status)
	statusCell.AppendChild(statusSelectElement)
	statusSelectElement.AddEventListener("change", true, func(e dom.Event) {
		jobId := strings.Split(statusSelectElement.ID(), "-")[1]
		job := jobs.NewJob("", "", job.Status, "", "")
		updateJob(jobId, job)
	})

	// Customer
	customerCell := row.InsertCell(4)
	customerCell.SetContentEditable("true")
	customerSelectElement := document.CreateElement("select").(*dom.HTMLSelectElement)
	customerSelectElement.Class().Add("form-control")
	customerSelectElement.SetID(fmt.Sprintf("customerDropdown-%s", job.ID))
	populateCustomerDropdownOptions(document, customerSelectElement, job.Customer)
	customerCell.AppendChild(customerSelectElement)
	customerSelectElement.AddEventListener("change", true, func(e dom.Event) {
		jobId := strings.Split(customerSelectElement.ID(), "-")[1]
		job := jobs.NewJob("", "", "", job.Customer, "")
		updateJob(jobId, job)
	})

	// Description
	decodedDescription, err := base64.StdEncoding.DecodeString(job.Description)
	if err != nil {
		log.Fatal(err)
	}
	descriptionCell := row.InsertCell(5)
	descriptionCell.SetContentEditable("true")
	descriptionTextArea := document.CreateElement("textarea").(*dom.HTMLTextAreaElement)
	descriptionTextArea.SetID(fmt.Sprintf("descriptionText-%s", job.ID))
	descriptionTextArea.Class().Add("form-control")
	descriptionCell.AppendChild(descriptionTextArea)
	descriptionTextArea.SetTextContent(string(decodedDescription))
	descriptionTextArea.AddEventListener("change", true, func(e dom.Event) {
		jobId := strings.Split(descriptionTextArea.ID(), "-")[1]
		newDescription := descriptionTextArea.Value
		job := jobs.NewJob("", "", "", "", newDescription)
		updateJob(jobId, job)
	})
}

func showUserInput(document dom.Document) {
	userInput := document.GetElementByID("userInput")
	jobsContainer := document.GetElementByID("jobsContainer")
	addRowContainer := document.GetElementByID("addRowBtnContainer")
	userInput.Class().Remove("d-none")
	jobsContainer.Class().Add("d-none")
	addRowContainer.Class().Add("d-none")
}

func submitJob(document dom.Document) {
	orderDate := document.GetElementByID("orderDateInput").(*dom.HTMLInputElement)
	deadlineDate := document.GetElementByID("deadlineInput").(*dom.HTMLInputElement)
	statusDropdown := document.GetElementByID("statusDropdown").(*dom.HTMLSelectElement)
	statusElement := statusDropdown.Options()[statusDropdown.SelectedIndex]
	customerDropdown := document.GetElementByID("customerDropdown").(*dom.HTMLSelectElement)
	customerElement := customerDropdown.Options()[customerDropdown.SelectedIndex]
	description := document.GetElementByID("descriptionInput").(*dom.HTMLTextAreaElement)
	job := jobs.NewJob(orderDate.Value, deadlineDate.Value, statusElement.Text, customerElement.Text,
		description.Value)
	payload, err := json.Marshal(job)
	if err != nil {
		log.Fatalf("PostJob:%v", err)
		return
	}

	// JavaScript callbacks cannot be blocking
	go func() {
		resp, err := http.Post("/jobs", "application/json", bytes.NewBuffer(payload))
		if err != nil || resp.StatusCode != http.StatusCreated {
			log.Fatalf("PostJob:%v\n", err)
		}
	}()

	hideUserInput(document)
	populateAllJobs(document)
}

func updateJob(id string, job *jobs.Job) {
	payload, err := json.Marshal(job)
	if err != nil {
		log.Fatalf("UpdateJob Marshal Error:%v", err)
		return
	}
	go func(id string, payload []byte) {
		url := fmt.Sprintf("/jobs/%s", id)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Fatalf("UpdateJob Request Error:%v\n", err)
		}
	}(id, payload)
}

func hideUserInput(document dom.Document) {
	userInput := document.GetElementByID("userInput")
	jobsContainer := document.GetElementByID("jobsContainer")
	addRowContainer := document.GetElementByID("addRowBtnContainer")
	userInput.Class().Add("d-none")
	jobsContainer.Class().Remove("d-none")
	addRowContainer.Class().Remove("d-none")

	populateAllJobs(document)
	orderDate := document.GetElementByID("orderDateInput").(*dom.HTMLInputElement)
	orderDate.Value = ""
	deadlineDate := document.GetElementByID("deadlineInput").(*dom.HTMLInputElement)
	deadlineDate.Value = ""
	statusDropdown := document.GetElementByID("statusDropdown").(*dom.HTMLSelectElement)
	statusDropdown.SelectedIndex = 0
	customerDropdown := document.GetElementByID("customerDropdown").(*dom.HTMLSelectElement)
	customerDropdown.SelectedIndex = 0
	description := document.GetElementByID("descriptionInput").(*dom.HTMLTextAreaElement)
	description.Value = ""
}
