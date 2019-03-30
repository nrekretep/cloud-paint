package cloudfoundry

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/nrekretep/cloudpaint/adapter/cloudfoundry/v3"
	"net/http"
	"net/url"
	"strings"
)

// CloudControllerConfig provides Cloud Foundry specific configuration options
type CloudControllerConfig struct {
	Username     string
	Password     string
	APIURLString string
	APIURL       *url.URL
}

// CloudController provides access to the cc API.
type CloudController struct {
	Config             *CloudControllerConfig
	APIUrl             *url.URL
	httpClient         *http.Client
	AccessToken        *AccessTokenInfo
	StackMap           *map[string]*StackInfo
	BuildpackMap       *map[string]*BuildpackInfo
	QuotaDefinitionMap *map[string]*QuotaDefinitionInfo
	OrganizationMap    *map[string]*OrganizationInfo
	SpaceMap           *map[string]*SpaceInfo
	AppMap             *map[string]*AppInfo
}

// NewCloudController returns a new CloudController client for the given url.
func NewCloudController(config CloudControllerConfig) (*CloudController, error) {

	err := checkConfig(&config)

	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyFromEnvironment,
	}

	httpClient := &http.Client{Transport: tr}

	c := &CloudController{httpClient: httpClient, APIUrl: config.APIURL, Config: &config}

	return c, nil
}

func checkConfig(c *CloudControllerConfig) error {

	if c == nil {
		return errors.New("config cannot be empty")
	}

	if c.Username == "" {
		return errors.New("username cannot be empty")
	}

	if c.Password == "" {
		return errors.New("password cannot be empty")
	}

	if c.APIURLString == "" {
		return errors.New("apiUrl cannot be empty")
	}

	u, err := url.Parse(c.APIURLString)
	if err != nil {
		return err
	}
	c.APIURL = u

	return nil
}

// Login to CC API and retrieve the access token.
func (c *CloudController) Login() error {

	parameters := url.Values{}
	parameters.Set("username", c.Config.Username)
	parameters.Set("password", c.Config.Password)
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

// GetV3App
func (c *CloudController) GetV3App(appID string) (*v3.App, error) {
	apiURLRelative := &url.URL{Path: "/v3/apps/" + appID}
	apiURL := c.APIUrl.ResolveReference(apiURLRelative)

	req, err := http.NewRequest("GET", apiURL.String(), nil)
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

	var a v3.App
	err = json.NewDecoder(resp.Body).Decode(&a)

	if err != nil {
		return nil, err
	}

	return &a, nil
}
