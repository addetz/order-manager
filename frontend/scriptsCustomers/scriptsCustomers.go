package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	customers "github.com/addetz/order-manager/services"
	"honnef.co/go/js/dom"
)

const DIVIDER = "#"

func main() {
	document := dom.GetWindow().Document()

	addCustomerBtn := document.GetElementByID("addCustomerBtn")
	submitBtn := document.GetElementByID("submitBtn")
	cancelBtn := document.GetElementByID("cancelBtn")
	addCustomerBtn.AddEventListener("click", true, func(e dom.Event) {
		container := document.GetElementByID("customerInput")
		tableContainer := document.GetElementByID("customerContainer")
		addCustomerBtnContainer := document.GetElementByID("addCustomerBtnContainer")
		container.Class().Remove("d-none")
		tableContainer.Class().Add("d-none")
		addCustomerBtnContainer.Class().Add("d-none")
	})

	submitBtn.AddEventListener("click", true, func(e dom.Event) {
		submitCustomer(document)
	})

	cancelBtn.AddEventListener("click", true, func(e dom.Event) {
		hideUserInput(document)
	})

	populateAllCustomers(document)
}

func submitCustomer(document dom.Document) {
	customerNameInput := document.GetElementByID("customerNameInput").(*dom.HTMLInputElement)
	customerNote := document.GetElementByID("customerNote").(*dom.HTMLTextAreaElement)

	customer := customers.NewCustomer(customerNameInput.Value, customerNote.Value)
	payload, err := json.Marshal(customer)
	if err != nil {
		log.Fatalf("PostCustomer:%v", err)
		return
	}

	go func() {
		resp, err := http.Post("/customers", "application/json", bytes.NewBuffer(payload))
		if err != nil || resp.StatusCode != http.StatusCreated {
			log.Fatalf("PostCustomer:%v\n", err)
		}
	}()

	populateAllCustomers(document)
	hideUserInput(document)
}

func hideUserInput(document dom.Document) {
	container := document.GetElementByID("customerInput")
	container.Class().Add("d-none")
	customerNameInput := document.GetElementByID("customerNameInput").(*dom.HTMLInputElement)
	customerNameInput.Value = ""
	customerNote := document.GetElementByID("customerNote").(*dom.HTMLTextAreaElement)
	customerNote.Value = ""
	tableContainer := document.GetElementByID("customerContainer")
	tableContainer.Class().Remove("d-none")
	addCustomerBtnContainer := document.GetElementByID("addCustomerBtnContainer")
	addCustomerBtnContainer.Class().Remove("d-none")
}

func populateAllCustomers(document dom.Document) {
	go func(callback func(document dom.Document, customers []*customers.Customer)) {
		resp, err := http.Get("/customers")
		if err != nil {
			log.Fatal(err)
		}
		customers, err := customers.NewCustomersResponse(resp)
		if err != nil {
			log.Fatal(err)
		}
		callback(document, customers)
	}(populateCustomersCallback)
}

func populateCustomersCallback(document dom.Document, customers []*customers.Customer) {
	newBody := document.CreateElement("tbody")
	ts := newBody.(*dom.HTMLTableSectionElement)
	for _, e := range customers {
		populateCustomer(document, ts, e)
	}
	oldBody := document.GetElementByID("customersTable").GetElementsByTagName("tbody")[0]
	document.GetElementByID("customersTable").ReplaceChild(newBody, oldBody)
}

func populateCustomer(document dom.Document,
	tableSection *dom.HTMLTableSectionElement,
	customer *customers.Customer) {
	row := tableSection.InsertRow(0)
	row.SetID(createElementID("row", customer.ID))

	// Name
	nameCell := row.InsertCell(0)
	nameCell.SetContentEditable("true")
	nameInput := document.CreateElement("input").(*dom.HTMLInputElement)
	nameInput.Class().Add("form-control")
	nameInput.SetID(createElementID("customerName", customer.ID))
	nameCell.AppendChild(nameInput)
	nameInput.Value = customer.Name
	nameInput.AddEventListener("change", true, func(e dom.Event) {
		customerId := extractCustomerIDFromElement(nameInput.ID())
		newName := nameInput.Value
		customer := customers.NewCustomer(newName, "")
		updateCustomer(document, customerId, customer)
	})

	// Description
	decodedDescription, err := base64.StdEncoding.DecodeString(customer.Note)
	if err != nil {
		log.Fatal(err)
	}
	noteCell := row.InsertCell(1)
	noteCell.SetContentEditable("true")
	noteTextArea := document.CreateElement("textarea").(*dom.HTMLTextAreaElement)
	noteTextArea.SetID(createElementID("customerNote", customer.ID))
	noteTextArea.Class().Add("form-control")
	noteCell.AppendChild(noteTextArea)
	noteTextArea.SetTextContent(string(decodedDescription))
	noteTextArea.AddEventListener("change", true, func(e dom.Event) {
		customerId := extractCustomerIDFromElement(noteTextArea.ID())
		newNote := noteTextArea.Value
		customer := customers.NewCustomer("", newNote)
		updateCustomer(document, customerId, customer)
	})

	// Delete button
	actionCell := row.InsertCell(2)
	deleteBtn := document.CreateElement("button").(*dom.HTMLButtonElement)
	deleteBtn.SetID(createElementID("deleteBtn", customer.ID))
	deleteBtn.Class().Add("btn")
	deleteBtn.Class().Add("btn-info")
	deleteBtn.Class().Add("mt-2")
	deleteBtn.SetTextContent("Delete Row?")
	actionCell.AppendChild(deleteBtn)
	deleteBtn.AddEventListener("click", true, func(e dom.Event) {
		customerId := extractCustomerIDFromElement(deleteBtn.ID())
		answer := dom.GetWindow().Confirm("Are you sure you want to delete row?")
		if answer {
			deleteCustomer(customerId, document)
		}
	})
}

func createElementID(prefix, id string) string {
	return fmt.Sprintf("%s%s%s", prefix, DIVIDER, id)
}

func extractCustomerIDFromElement(id string) string {
	return strings.Split(id, DIVIDER)[1]
}

func updateCustomer(document dom.Document, id string, customer *customers.Customer) {
	payload, err := json.Marshal(customer)
	if err != nil {
		log.Fatalf("UpdateCustomer Marshal Error:%v", err)
		return
	}
	go func(id string, payload []byte) {
		url := fmt.Sprintf("/customers/%s", id)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Fatalf("UpdateCustomer Request Error:%v\n", err)
		}
		populateAllCustomers(document)
	}(id, payload)
}

func deleteCustomer(id string, document dom.Document) {
	go func(id string) {
		url := fmt.Sprintf("/customers/%s", id)
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			log.Fatalf("DeleteCustomer Request Error:%v\n", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Fatalf("DeleteCustomer Error:%v\n", err)
		}
		populateAllCustomers(document)
	}(id)
}
