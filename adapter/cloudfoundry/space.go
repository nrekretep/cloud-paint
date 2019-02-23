package cloudfoundry

import (
	"encoding/json"
	"fmt"
)

// SpaceEntity
type SpaceEntity struct {
	Name                     string `json:"name"`
	OrganizationGUID         string `json:"organization_guid"`
	OrganizationURL          string `json:"organization_url"`
	Organization             *OrganizationInfo
	SpaceQuotaDefinitionGUID string `json:"space_quota_definition_guid"`
	SpaceQuotaDefinition     *QuotaDefinitionInfo
	AllowSSH                 bool   `json:"allow_ssh"`
	DevelopersURL            string `json:"developers_url"`
	ManagersURL              string `json:"managers_url"`
	//Managers *map[string]*UserInfo
	AuditorsURL string `json:"auditors_url"`
	//Auditors *map[string]*UserInfo
	DomainsURL string `json:"domains_url"`
	//Domains                  *map[string]*DomainInfo
	AppEventsURL             string `json:"app_events_url"`
	EventsURL                string `json:"events_url"`
	AppsURL                  string `json:"apps_url"`
	RoutesURL                string `json:"routes_url"`
	ServiceInstancesURL      string `json:"service_instances_url"`
	SecurityGroupsURL        string `json:"security_groups_url"`
	StagingSecurityGroupsURL string `json:"staging_security_groups_url"`
}

// SpaceInfo
type SpaceInfo struct {
	Metadata Metadata
	Entity   SpaceEntity
}

// GetSpaces
func (c *CloudController) GetSpaces() error {

	spaceResources, err := c.GetResourceList("/v2/spaces")

	resultMap := make(map[string]*SpaceInfo)

	for _, value := range *spaceResources {
		s := new(SpaceInfo)
		s.Metadata = value.Metadata
		err = json.Unmarshal(value.Entity, &s.Entity)

		if err != nil {
			fmt.Println(err)
			return err
		}

		resultMap[s.Metadata.GUID] = s
		fmt.Println(s.Metadata.GUID)
		fmt.Println(s.Metadata.CreatedAt)
		fmt.Println(s.Metadata.UpdatedAt)
		fmt.Println(s.Entity.Name)
		fmt.Println(s.Entity.OrganizationGUID)
		fmt.Println("============")
	}

	c.SpaceMap = &resultMap
	return err

}
