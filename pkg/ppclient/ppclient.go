package ppclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type PpClient interface {
	GetReleases(productShortname string) (PpReleaseList, error)
}

type client struct {
	client   *http.Client
	endpoint string
}

func NewPpClient(endpoint string) PpClient {
	return &client{
		client: &http.Client{
			Timeout: time.Duration(10 * time.Second),
		},
		endpoint: endpoint,
	}
}

func (c *client) GetReleases(productShortname string) (PpReleaseList, error) {
	var releases []struct {
		Shortname string `json:"shortname,omitempty"`
		Phase     string `json:"phase_display,omitempty"`
	}

	endpointPath := "/releases?product__shortname=" + productShortname + "&fields=shortname,phase_display"
	request, err := http.NewRequest(http.MethodGet, c.endpoint+endpointPath, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %v", err)
	}
	resp, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error in http request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("error: http response status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading http response: %v", err)
	}
	err = json.Unmarshal(body, &releases)
	if err != nil {
		return nil, fmt.Errorf("error decoding http response to json: %v", err)
	}

	var ppReleases PpReleaseList

	for _, r := range releases {
		ppReleases = append(ppReleases, NewPpRelease(r.Shortname, r.Phase))
	}

	return ppReleases, nil
}
