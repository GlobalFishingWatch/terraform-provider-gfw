package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const DATASET_PATH = "datasets"

func (c *GFWClient) GetDatasets() (*[]Dataset, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s?includes[0]=BACKEND_CONFIGURATION&cache=false", c.HostURL, DATASET_PATH), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	datasets := Pagination[Dataset]{}
	err = json.Unmarshal(body, &datasets)
	if err != nil {
		return nil, err
	}

	return &datasets.Entries, nil
}

func (c *GFWClient) GetDataset(id string) (*Dataset, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s?includes[0]=BACKEND_CONFIGURATION&cache=false", c.HostURL, DATASET_PATH, id), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	dataset := Dataset{}
	err = json.Unmarshal(body, &dataset)
	if err != nil {
		return nil, err
	}

	return &dataset, nil
}

func (c *GFWClient) DeleteDataset(id string) (*Dataset, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s?includes[0]=BACKEND_CONFIGURATION", c.HostURL, DATASET_PATH, id), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, nil
	}
	dataset := Dataset{}
	err = json.Unmarshal(body, &dataset)
	if err != nil {
		return nil, err
	}

	return &dataset, nil
}

func (c *GFWClient) UpdateDataset(id string, dataset CreateDataset) error {

	bodyReq, err := json.Marshal(dataset)
	if err != nil {
		return err
	}
	fmt.Println(string(bodyReq))
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s/%s?includes[0]=BACKEND_CONFIGURATION", c.HostURL, DATASET_PATH, id), strings.NewReader(string(bodyReq)))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return err
	}
	_, err = c.doRequest(req)
	return err

}

func (c *GFWClient) CreateDataset(dataset CreateDataset) (*Dataset, error) {
	exists, err := c.checkExistDataset(dataset.ID)
	if err != nil {
		return nil, err
	}
	if exists != nil {
		return exists, nil
	}

	bodyReq, err := json.Marshal(dataset)
	if err != nil {
		return nil, err
	}
	data := string(bodyReq)
	fmt.Println(data)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.HostURL, DATASET_PATH), strings.NewReader(string(bodyReq)))
	req.Header.Add("content-type", "application/json")
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newDataset := Dataset{}
	err = json.Unmarshal(body, &newDataset)
	if err != nil {
		return nil, err
	}

	return &newDataset, nil
}

func (c *GFWClient) checkExistDataset(id string) (*Dataset, error) {
	dataset, err := c.GetDataset(id)
	if err != nil {
		if re, ok := err.(AppError); ok {
			if re.Code == NotFoundCode {
				return nil, nil
			}
		}
		return nil, err
	}

	return dataset, nil
}
