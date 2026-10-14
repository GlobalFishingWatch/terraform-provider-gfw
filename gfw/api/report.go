package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const REPORT_PATH = "reports"

func (c *GFWClient) GetReport(id string) (*Report, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", c.HostURL, REPORT_PATH, id), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	report := Report{}
	err = json.Unmarshal(body, &report)
	if err != nil {
		return nil, err
	}
	report.Public = strings.HasSuffix(report.ID, "-public")
	return &report, nil
}

func (c *GFWClient) DeleteReport(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", c.HostURL, REPORT_PATH, id), nil)
	if err != nil {
		return err
	}
	_, err = c.doRequest(req)
	return err
}

func (c *GFWClient) UpdateReport(id string, report UpdateReport) error {
	bodyReq, err := json.Marshal(report)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%s/%s", c.HostURL, REPORT_PATH, id), strings.NewReader(string(bodyReq)))
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	_, err = c.doRequest(req)
	return err
}

func (c *GFWClient) CreateReport(report CreateReport) (*Report, error) {
	// Adopt an existing report with the same explicit id instead of failing on duplicate
	if report.ID != "" {
		exists, err := c.GetReport(report.ID)
		if err == nil {
			return exists, nil
		}
		if re, ok := err.(AppError); !ok || re.Code != NotFoundCode {
			return nil, err
		}
	}

	bodyReq, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.HostURL, REPORT_PATH), strings.NewReader(string(bodyReq)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newReport := Report{}
	err = json.Unmarshal(body, &newReport)
	if err != nil {
		return nil, err
	}

	return &newReport, nil
}
