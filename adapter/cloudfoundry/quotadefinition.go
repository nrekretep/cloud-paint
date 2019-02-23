package cloudfoundry

import (
	"encoding/json"
	"fmt"
)

// QuotaDefinitionEntity - Entity data for Org Quota Definitions.
type QuotaDefinitionEntity struct {
	Name                    string `json:"name"`
	NonBasicServicesAllowed bool   `json:"non_basic_services_allowed"`
	TotalService            int    `json:"total_services"`
	TotalRoutes             int    `json:"total_routes"`
	TotalPrivateDomains     int    `json:"total_private_domains"`
	MemoryLimit             int    `json:"memory_limit"`
	TrialDBAllowed          bool   `json:"trial_db_allowed"`
	InstanceMemoryLimit     int    `json:"instance_memory_limit"`
	AppInstanceLimit        int    `json:"app_instance_limit"`
	AppTaskLimit            int    `json:"app_task_limit"`
	TotalServiceKeys        int    `json:"total_service_keys"`
	TotalReservedRoutePorts int    `json:"total_reserved_route_ports"`
}

// QuotaDefinitionInfo - Details about Org Quota Definitions.
type QuotaDefinitionInfo struct {
	Metadata Metadata
	Entity   QuotaDefinitionEntity
}

// GetQuotaDefinitions
func (c *CloudController) GetQuotaDefinitions() error {

	quotaDefinitionResources, err := c.GetResourceList("/v2/quota_definitions")

	resultMap := make(map[string]*QuotaDefinitionInfo)

	for _, value := range *quotaDefinitionResources {
		oqdi := new(QuotaDefinitionInfo)
		oqdi.Metadata = value.Metadata
		err = json.Unmarshal(value.Entity, &oqdi.Entity)

		if err != nil {
			fmt.Println(err)
			return err
		}

		resultMap[oqdi.Metadata.GUID] = oqdi
		fmt.Println(oqdi.Metadata.GUID)
		fmt.Println(oqdi.Metadata.CreatedAt)
		fmt.Println(oqdi.Metadata.UpdatedAt)
		fmt.Println(oqdi.Entity.Name)
		fmt.Println(oqdi.Entity.TotalService)
		fmt.Println(oqdi.Entity.TotalRoutes)
		fmt.Println(oqdi.Entity.InstanceMemoryLimit)
		fmt.Println("============")
	}

	c.QuotaDefinitionMap = &resultMap
	return err

}
