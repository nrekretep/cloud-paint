package cloudfoundry

import (
	"encoding/json"
	//"fmt"
	"net/http"
	"net/url"
)

// ResourceList - List of generic resources
type ResourceList struct {
	TotalResults int    `json:"total_results"`
	TotalPages   int    `json:"total_pages"`
	PreviousURL  string `json:"prev_url"`
	NextURL      string `json:"next_url"`
	Resources    []Resource
}

// Resource - Generic
type Resource struct {
	Metadata Metadata
	Entity   json.RawMessage
}

// GetResourceList returns a map of Resources
func (c *CloudController) GetResourceList(apiPath string) (*map[string]Resource, error) {
	apiURLRelative := &url.URL{Path: apiPath}
	apiURL := c.APIUrl.ResolveReference(apiURLRelative)

	req, err := http.NewRequest("GET", apiURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AccessToken.AccessToken)

	q := req.URL.Query()
	q.Add("results-per-page", "100")
	req.URL.RawQuery = q.Encode()

	var i ResourceList
	resourceList := make(map[string]Resource)
	hasNext := true

	for hasNext {
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&i)

		for _, value := range i.Resources {
			resourceList[value.Metadata.GUID] = value
		}

		hasNext = !(i.NextURL == "")
		if hasNext {
			apiURL, _ = url.Parse(c.APIUrl.String() + i.NextURL)
			req.URL = apiURL
			i.NextURL = ""
		}
	}

	return &resourceList, err
}
