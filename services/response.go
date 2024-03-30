package jobs

import (
	"encoding/json"
	"io"
	"net/http"
)

func NewBackendResponse(resp *http.Response) ([]*Job, error) {
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

func NewJob(orderDate string, deadline string, status string,
	customer string, description string) *Job {
	j := &Job{}
	if orderDate != "" {
		j.OrderDate = *getFormattedDate(orderDate)
	}
	if deadline != "" {
		j.DeadlineDate = *getFormattedDate(deadline)
	}
	j.Status = status
	j.Customer = customer
	j.Description = description
	return j
}
