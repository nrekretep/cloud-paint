package cloudfoundry

import (
	"encoding/json"
	"fmt"
)

// AppEntity
type AppEntity struct {
	Name                     string `json:"name"`
	Production               bool   `json:"production"`
	SpaceGUID                string `json:"space_guid"`
	SpaceURL                 string `json:"space_url"`
	Space                    *SpaceInfo
	StackGUID                string `json:"stack_guid"`
	Stack                    *StackInfo
	Buildpack                string          `json:"buildpack"`
	DetectedBuildpack        string          `json:"detected_buildpack"`
	DetectedBuildpackGUID    string          `json:"detected_buildpack_guid"`
	EnvironmentJSON          json.RawMessage `json:"environment_json"`
	Memory                   int             `json:"memory"`
	Instances                int             `json:"instances"`
	DiskQuota                int             `json:"disk_quota"`
	State                    string          `json:"state"`
	Version                  string          `json:"version"`
	Comman                   string          `json:"command"`
	Console                  bool            `json:"console"`
	Debug                    string          `json:"debug"`
	StaginTaskID             string          `json:"staging_task_id"`
	PackageStage             string          `json:"package_state"`
	HealthCheckType          string          `json:"health_check_type"`
	HealthCheckTimeout       int             `json:"health_check_timeout"`
	StagingFailedReason      string          `string:"staging_failed_reason"`
	StagingFailedDescription string          `json:"staging_failed_description"`
	Diego                    bool            `json:"diego"`
	DockerImage              string          `json:"docker_image"`
	PackageUpdatedAt         string          `json:"package_updated_at"`
	DetectedStartCommand     string          `json:"detected_start_command"`
	EnableSSH                bool            `json:"enable_ssh"`
	Ports                    []int           `json:"ports"`
	RoutesURL                string          `json:"routes_url"`
	EventsURL                string          `json:"events_url"`
	ServiceBindingURL        string          `json:"service_bindings_url"`
	RouteMappingURL          string          `json:"route_mappings_url"`
}

// AppInfo
type AppInfo struct {
	Metadata Metadata
	Entity   AppEntity
}

// GetApps
func (c *CloudController) GetApps() error {

	appResources, err := c.GetResourceList("/v2/apps")

	resultMap := make(map[string]*AppInfo)

	for _, value := range *appResources {
		a := new(AppInfo)
		a.Metadata = value.Metadata
		err = json.Unmarshal(value.Entity, &a.Entity)

		if err != nil {
			fmt.Println(err)
			return err
		}

		resultMap[a.Metadata.GUID] = a
		// fmt.Println(a.Metadata.GUID)
		// fmt.Println(a.Metadata.CreatedAt)
		// fmt.Println(a.Metadata.UpdatedAt)
		// fmt.Println(a.Entity.Name)
		// fmt.Println("============")
	}

	c.AppMap = &resultMap
	return err

}
