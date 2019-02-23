package cloudfoundry

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// CloudController provides access to the cc API.
type CloudController struct {
	APIUrl             *url.URL
	httpClient         *http.Client
	AccessToken        *AccessTokenInfo
	StackMap           *map[string]*StackInfo
	BuildpackMap       *map[string]*BuildpackInfo
	QuotaDefinitionMap *map[string]*QuotaDefinitionInfo
	OrganizationMap    *map[string]*OrganizationInfo
	SpaceMap           *map[string]*SpaceInfo
}

// NewCloudController returns a new CloudController client for the given url.
func NewCloudController(apiURL string) *CloudController {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyFromEnvironment,
	}

	httpClient := &http.Client{Transport: tr}

	u, err := url.Parse(apiURL)
	if err != nil {
		return nil
	}

	c := &CloudController{httpClient: httpClient, APIUrl: u}

	return c
}

// Login to CC API and retrieve the access token.
func (c *CloudController) Login(usename string, password string) error {

	parameters := url.Values{}
	parameters.Set("username", usename)
	parameters.Set("password", password)
	parameters.Set("scope", "")
	parameters.Set("grant_type", "password")

	info, err := c.GetV2Info()
	if err != nil {
		return nil
	}

	authURLRelative := &url.URL{Path: "/oauth/token"}
	authURL, err := url.Parse(info.AuthorizationEndpoint)
	if err != nil {
		return nil
	}
	authTokenURL := authURL.ResolveReference(authURLRelative)
	req, err := http.NewRequest("POST", authTokenURL.String(), strings.NewReader(parameters.Encode()))
	if err != nil {
		return nil
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", info.AuthorizationEndpoint)
	req.SetBasicAuth("cf", "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var ati AccessTokenInfo
	err = json.NewDecoder(resp.Body).Decode(&ati)
	if err != nil {
		return err
	}

	c.AccessToken = &ati
	return nil
}

// GetV2Info gets the general API info from the select cc API.
func (c *CloudController) GetV2Info() (*V2Info, error) {
	infoURLRelative := &url.URL{Path: "/v2/info"}
	infoURL := c.APIUrl.ResolveReference(infoURLRelative)

	req, err := http.NewRequest("GET", infoURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var i V2Info
	err = json.NewDecoder(resp.Body).Decode(&i)

	return &i, err
}

// GetAppInfo delivery details for the given app guid.
func (c *CloudController) GetAppInfo(guid string) (*AppInfo, error) {

	if guid == "" {
		return nil, errors.New("app guid is empty")
	}

	appURL := c.APIUrl.ResolveReference(&url.URL{Path: "/v2/apps/" + guid})

	req, err := http.NewRequest("GET", appURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AccessToken.AccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ai AppInfo
	err = json.NewDecoder(resp.Body).Decode(&ai)

	return &ai, err

}

// V2Info represents general info about V2 cloud controller API.
type V2Info struct {
	Name                     string `json:"name"`
	Build                    string `json:"build"`
	Support                  string `json:"support"`
	Version                  int    `json:"version"`
	Description              string `json:"description"`
	AuthorizationEndpoint    string `json:"authorization_endpoint"`
	TokenEndpoint            string `json:"token_endpoint"`
	MinCliVersion            string `json:"min_cli_version"`
	MinRecommendedCliVersion string `json:"min_recommended_cli_version"`
	AppSSHEndpoint           string `json:"app_ssh_endpoint"`
	AppSSHHostKeyFingerprint string `json:"app_ssh_host_key_fingerprint"`
	DopplerLoggingEndpoint   string `json:"doppler_logging_endpoint"`
	APIVersion               string `json:"api_version"`
	OsbapiVersion            string `json:"osbapi_version"`
	RoutingEndpoint          string `json:"routing_endpoint"`
}

// AccessTokenInfo contains details for the current access token.
type AccessTokenInfo struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	JTI          string `json:"jti"`
}

// AppInfo contains application specific details
type AppInfo struct {
	Metadata Metadata
	Entity   AppEntity
}

// AppEntity contains Entity data for the application.
type AppEntity struct {
	Name string `json:"name"`
	// 	  "production": false,
	SpaceGUID string `json:"space_guid"`
	StackGUID string `json:"stack_guid"`
	// 	  "buildpack": null,
	// 	  "detected_buildpack": null,
	// 	  "detected_buildpack_guid": null,
	// 	  "environment_json": null,
	// 	  "memory": 1024,
	// 	  "instances": 1,
	// 	  "disk_quota": 1024,
	// 	  "state": "STOPPED",
	// 	  "version": "df19a7ea-2003-4ecb-a909-e630e43f2719",
	// 	  "command": null,
	// 	  "console": false,
	// 	  "debug": null,
	// 	  "staging_task_id": null,
	// 	  "package_state": "PENDING",
	// 	  "health_check_http_endpoint": "",
	HealthCheckTyep    string `json:"health_check_type"`
	HealthCheckTimeout string `json:"health_check_timeout"`
	// 	  "staging_failed_reason": null,
	// 	  "staging_failed_description": null,
	// 	  "diego": false,
	// 	  "docker_image": null,
	// 	  "docker_credentials": {
	// 		"username": null,
	// 		"password": null
	// 	  },
	// 	  "package_updated_at": "2016-06-08T16:41:45Z",
	// 	  "detected_start_command": "",
	// 	  "enable_ssh": true,
	// 	  "ports": null,
	// 	  "space_url": "/v2/spaces/7846301e-c84c-4ba9-9c6a-2dfdae948d52",
	// 	  "stack_url": "/v2/stacks/7e03186d-a438-4285-b3b7-c426532e1df2",
	// 	  "routes_url": "/v2/apps/15b3885d-0351-4b9b-8697-86641668c123/routes",
	// 	  "events_url": "/v2/apps/15b3885d-0351-4b9b-8697-86641668c123/events",
	// 	  "service_bindings_url": "/v2/apps/15b3885d-0351-4b9b-8697-86641668c123/service_bindings",
	// 	  "route_mappings_url": "/v2/apps/15b3885d-0351-4b9b-8697-86641668c123/route_mappings"
	// 	}
	//   }

}
