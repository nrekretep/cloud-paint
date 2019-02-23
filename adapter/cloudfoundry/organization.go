package cloudfoundry

import (
	"encoding/json"
	"fmt"
)

// OrganizationEntity
type OrganizationEntity struct {
	Name                string `json:"name"`
	BillingEnabled      bool   `json:"billing_enabled"`
	QuotaDefinitionGUID string `json:"quota_definition_guid"`
	QuotaDefinition     *QuotaDefinitionInfo
	Status              string `json:"status"`
	QuotaDefinitionURL  string `json:"quota_definition_url"`
	SpacesURL           string `json:"spaces_url"`
	Spaces              *map[string]*SpaceInfo
	DomainsURL          string `json:"domains_url"`
	//Domains                  *map[string]*DomainInfo
	PrivateDomainsURL string `json:"private_domains_url"`
	//PrivateDomains                  *map[string]*PrivateDomainInfo
	UsersURL string `json:"users_url"`
	//Users *map[string]*UserInfo
	ManagersURL string `json:"managers_url"`
	//Managers *map[string]*UserInfo
	BillingManagersURL string `json:"billing_managers_url"`
	//BillingManagers *map[string]*UserInfo
	AuditorsURL string `json:"auditors_url"`
	//Auditors *map[string]*UserInfo
	AppEventsURL             string `json:"app_events_url"`
	SpaceQuotaDefinitionsURL string `json:"space_quota_definitions_url"`
	//SpaceQuotas *map[string]*SpaceQuotaInfo
}

// OrganizationInfo
type OrganizationInfo struct {
	Metadata Metadata
	Entity   OrganizationEntity
}

// GetOrganizations
func (c *CloudController) GetOrganizations() error {

	organizationResources, err := c.GetResourceList("/v2/organizations")

	resultMap := make(map[string]*OrganizationInfo)

	for _, value := range *organizationResources {
		o := new(OrganizationInfo)
		o.Metadata = value.Metadata
		err = json.Unmarshal(value.Entity, &o.Entity)

		if err != nil {
			fmt.Println(err)
			return err
		}

		resultMap[o.Metadata.GUID] = o
		fmt.Println(o.Metadata.GUID)
		fmt.Println(o.Metadata.CreatedAt)
		fmt.Println(o.Metadata.UpdatedAt)
		fmt.Println(o.Entity.Name)
		fmt.Println(o.Entity.Status)
		fmt.Println(o.Entity.ManagersURL)
		fmt.Println(o.Entity.SpaceQuotaDefinitionsURL)
		fmt.Println("============")
	}

	c.OrganizationMap = &resultMap
	return err

}
