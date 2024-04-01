package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	jobs "github.com/addetz/order-manager/services"
	"honnef.co/go/js/dom"
)

const DIVIDER = "#"

func main() {
	document := dom.GetWindow().Document()
	populateAllJobs(document, "")

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
	addCustomerFilter(document)
}

func addCustomerFilter(document dom.Document) {
	filterCustomerDropdown := document.GetElementByID("filterCustomerDropdown").(*dom.HTMLSelectElement)
	populateCustomerDropdownOptions(document, filterCustomerDropdown, "")
	filterCustomerDropdown.AddEventListener("change", true, func(e dom.Event) {
		filter := filterCustomerDropdown.SelectedOptions()[0].Value
		go func(document dom.Document, filter string) {
			if filter == "All" {
				populateAllJobs(document, "")
				return
			}
			if filter == "Unknown" {
				populateAllJobs(document, "unknown")
				return
			}
			resp, err := http.Get(fmt.Sprintf("/customers/search?name=%s", filter))
			if err != nil {
				log.Fatal(err)
			}
			customer, err := jobs.NewCustomerSearchResponse(resp)
			if err != nil {
				log.Fatal(err)
			}
			populateAllJobs(document, customer.ID)
		}(document, filter)
	})
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
	log.Println(currentValue)
	o := document.CreateElement("option")
	o.SetTextContent("Unknown")
	customerDropdown.AppendChild(o)
	customerDropdown.SelectedIndex = 0

	for i, c := range customers {
		o := document.CreateElement("option")
		o.SetTextContent(c.Name)
		customerDropdown.AppendChild(o)
		if c.ID == currentValue {
			customerDropdown.SelectedIndex = i + 1
		}
	}
}

func populateAllJobs(document dom.Document, filter string) {
	log.Println("populate all jobs invoked with filter ", filter)
	go func(callback func(document dom.Document, jobs []*jobs.Job)) {
		url := "/jobs"
		if filter != "" {
			url = fmt.Sprintf("/jobs?customerID=%s", filter)
		}
		resp, err := http.Get(url)
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

func populateJob(document dom.Document,
	tableSection *dom.HTMLTableSectionElement,
	job *jobs.Job) {
	row := tableSection.InsertRow(0)
	row.SetID(createElementID("row", job.ID))

	// Order Date
	orderDateCell := row.InsertCell(0)
	orderDateCell.SetContentEditable("true")
	orderDatePicker := document.CreateElement("input").(*dom.HTMLInputElement)
	orderDatePicker.SetAttribute("type", "date")
	orderDatePicker.Class().Add("form-control")
	orderDatePicker.SetID(createElementID("orderDate", job.ID))
	orderDateCell.AppendChild(orderDatePicker)
	orderDatePicker.Value = job.OrderDate.Format(jobs.JobsDateFormat)
	orderDatePicker.AddEventListener("change", true, func(e dom.Event) {
		jobId := extractJobIDFromElement(orderDatePicker.ID())
		newOrderDate := orderDatePicker.Value
		job := jobs.NewJob(newOrderDate, "", "", "", "")
		updateJob(document, jobId, job)
	})

	// Deadline Date
	deadlineDateCell := row.InsertCell(1)
	deadlineDateCell.SetContentEditable("true")
	deadlineDatePicker := document.CreateElement("input").(*dom.HTMLInputElement)
	deadlineDatePicker.Class().Add("form-control")
	deadlineDatePicker.SetID(createElementID("deadlineDate", job.ID))
	deadlineDateCell.AppendChild(deadlineDatePicker)
	deadlineDatePicker.Value = job.DeadlineDate.Format(jobs.JobsDateFormat)
	deadlineDatePicker.AddEventListener("change", true, func(e dom.Event) {
		jobId := extractJobIDFromElement(deadlineDatePicker.ID())
		newDeadlineDate := deadlineDatePicker.Value
		job := jobs.NewJob("", newDeadlineDate, "", "", "")
		updateJob(document, jobId, job)
	})

	// Status
	statusCell := row.InsertCell(2)
	statusCell.SetContentEditable("true")
	statusSelectElement := document.CreateElement("select").(*dom.HTMLSelectElement)
	statusSelectElement.Class().Add("form-control")
	statusSelectElement.SetID(createElementID("statusDropdown", job.ID))
	populateStatusDropdownOptions(document, statusSelectElement, job.Status)
	statusCell.AppendChild(statusSelectElement)
	statusSelectElement.AddEventListener("change", true, func(e dom.Event) {
		jobId := extractJobIDFromElement(statusSelectElement.ID())
		newStatus := statusSelectElement.SelectedOptions()[0].Value
		job := jobs.NewJob("", "", newStatus, "", "")
		updateJob(document, jobId, job)
	})

	// Customer
	customerCell := row.InsertCell(3)
	customerCell.SetContentEditable("true")
	customerSelectElement := document.CreateElement("select").(*dom.HTMLSelectElement)
	customerSelectElement.Class().Add("form-control")
	customerSelectElement.SetID(createElementID("customerDropdown", job.ID))
	populateCustomerDropdownOptions(document, customerSelectElement, job.CustomerID)
	customerCell.AppendChild(customerSelectElement)
	customerSelectElement.AddEventListener("change", true, func(e dom.Event) {
		jobId := extractJobIDFromElement(customerSelectElement.ID())
		newCustomer := customerSelectElement.SelectedOptions()[0].Value
		go func(document dom.Document, newCustomer string, jobId string) {
			newJob := jobs.NewJob("", "", "", "Unknown", "")
			if newCustomer != "Unknown" {
				resp, err := http.Get(fmt.Sprintf("/customers/search?name=%s", newCustomer))
				if err != nil {
					log.Fatal(err)
				}
				customer, err := jobs.NewCustomerSearchResponse(resp)
				if err != nil {
					log.Fatal(err)
				}
				newJob = jobs.NewJob("", "", "", customer.ID, "")
			}
			updateJob(document, jobId, newJob)
		}(document, newCustomer, jobId)
	})

	// Description
	decodedDescription, err := base64.StdEncoding.DecodeString(job.Description)
	if err != nil {
		log.Fatal(err)
	}
	descriptionCell := row.InsertCell(4)
	descriptionCell.SetContentEditable("true")
	descriptionTextArea := document.CreateElement("textarea").(*dom.HTMLTextAreaElement)
	descriptionTextArea.SetID(createElementID("descriptionText", job.ID))
	descriptionTextArea.Class().Add("form-control")
	descriptionCell.AppendChild(descriptionTextArea)
	descriptionTextArea.SetTextContent(string(decodedDescription))
	descriptionTextArea.AddEventListener("change", true, func(e dom.Event) {
		jobId := extractJobIDFromElement(descriptionTextArea.ID())
		newDescription := descriptionTextArea.Value
		job := jobs.NewJob("", "", "", "", newDescription)
		updateJob(document, jobId, job)
	})

	// Delete button
	actionCell := row.InsertCell(5)
	deleteBtn := document.CreateElement("button").(*dom.HTMLButtonElement)
	deleteBtn.SetID(createElementID("deleteBtn", job.ID))
	deleteBtn.Class().Add("btn")
	deleteBtn.Class().Add("btn-info")
	deleteBtn.Class().Add("mt-2")
	deleteBtn.SetTextContent("Delete Row?")
	actionCell.AppendChild(deleteBtn)
	deleteBtn.AddEventListener("click", true, func(e dom.Event) {
		jobId := extractJobIDFromElement(deleteBtn.ID())
		answer := dom.GetWindow().Confirm("Are you sure you want to delete row?")
		if answer {
			deleteJob(jobId, document)
		}
	})

	// At the end apply the style
	applyRowStyle(row, job)
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
	job := jobs.NewJob(orderDate.Value, deadlineDate.Value, statusElement.Text, "",
		description.Value)
	customerName := customerElement.Value

	go func(job *jobs.Job, customerName string) {
		if customerName != "Unknown" {
			resp, err := http.Get(fmt.Sprintf("/customers/search?name=%s", customerName))
			if err != nil {
				log.Fatal(err)
			}
			customer, err := jobs.NewCustomerSearchResponse(resp)
			if err != nil {
				log.Fatal(err)
			}
			job.CustomerID = customer.ID
		}
		payload, err := json.Marshal(job)
		if err != nil {
			log.Fatalf("PostJob:%v", err)
			return
		}

		resp, err := http.Post("/jobs", "application/json", bytes.NewBuffer(payload))
		if err != nil || resp.StatusCode != http.StatusCreated {
			log.Fatalf("PostJob:%v\n", err)
		}

		populateAllJobs(document, "")
	}(job, customerName)

	hideUserInput(document)
	filterCustomerDropdown := document.GetElementByID("filterCustomerDropdown").(*dom.HTMLSelectElement)
	filterCustomerDropdown.SelectedIndex = 0
}

func deleteJob(id string, document dom.Document) {
	go func(id string) {
		url := fmt.Sprintf("/jobs/%s", id)
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			log.Fatalf("DeleteJob Request Error:%v\n", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Fatalf("DeleteJob Error:%v\n", err)
		}
		populateAllJobs(document, "")
	}(id)
}

func updateJob(document dom.Document, id string, job *jobs.Job) {
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
		populateAllJobs(document, "")
	}(id, payload)
}

func hideUserInput(document dom.Document) {
	userInput := document.GetElementByID("userInput")
	jobsContainer := document.GetElementByID("jobsContainer")
	addRowContainer := document.GetElementByID("addRowBtnContainer")
	userInput.Class().Add("d-none")
	jobsContainer.Class().Remove("d-none")
	addRowContainer.Class().Remove("d-none")

	populateAllJobs(document, "")
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

func calculateWorkingDays(dateTo time.Time) int {
	nowString := time.Now().Format(jobs.JobsDateFormat)
	dateFrom := *jobs.GetFormattedDate(nowString)

	if dateTo.Before(dateFrom) {
		return -1
	}

	days := 0
	for {
		if dateFrom.Equal(dateTo) {
			days++
			break
		}
		if dateFrom.Weekday() != 6 && dateFrom.Weekday() != 0 {
			days++
		}
		dateFrom = dateFrom.Add(time.Hour * 24)
	}

	return days
}

func applyRowStyle(row *dom.HTMLTableRowElement, job *jobs.Job) {
	// this job is finished
	if job.Status == jobs.JobStatusList[2] {
		row.Class().Add("finished-row")
		return
	}

	// this job is shipped
	if job.Status == jobs.JobStatusList[1] {
		row.Class().Add("shipped-row")
		return
	}

	daysLeft := calculateWorkingDays(*job.DeadlineDate)

	// the job is new and it is overdue
	if job.Status == jobs.JobStatusList[0] && daysLeft == -1 {
		row.Class().Add("overdue-row")
		return
	}

	// the job is new and it is due next day
	if job.Status == jobs.JobStatusList[0] && daysLeft == 1 {
		row.Class().Add("danger-row")
		return
	}

	// the job is new and there is less than a week left
	if job.Status == jobs.JobStatusList[0] && daysLeft < 5 {
		row.Class().Add("warning-row")
		return
	}
}

func createElementID(prefix, id string) string {
	return fmt.Sprintf("%s%s%s", prefix, DIVIDER, id)
}

func extractJobIDFromElement(id string) string {
	return strings.Split(id, DIVIDER)[1]
}
