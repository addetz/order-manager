package jobs

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
)

func NewJobsResponse(resp *http.Response) ([]*Job, error) {
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var bs []*Job
	if err := json.Unmarshal(body, &bs); err != nil {
		return nil, err
	}

	return bs, nil
}

func NewCustomersResponse(resp *http.Response) ([]*Customer, error) {
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var bs []*Customer
	if err := json.Unmarshal(body, &bs); err != nil {
		return nil, err
	}

	return bs, nil
}

func NewJob(orderDate string, deadline string, status string,
	customer string, description string) *Job {
	j := &Job{}
	if orderDate != "" {
		j.OrderDate = GetFormattedDate(orderDate)
	}
	if deadline != "" {
		j.DeadlineDate = GetFormattedDate(deadline)
	}
	j.Status = status
	j.Customer = customer
	j.Description = base64.StdEncoding.EncodeToString([]byte(description))
	return j
}
