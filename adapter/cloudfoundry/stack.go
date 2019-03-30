package cloudfoundry

import (
	"encoding/json"
	//"fmt"
	// "errors"
	// "net/http"
	// "net/url"
)

// StackEntity - Entity data for stacks.
type StackEntity struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// StackInfo - Details about stacks.
type StackInfo struct {
	Metadata Metadata
	Entity   StackEntity
}

// GetStacks - Loads infos about all stacks
func (c *CloudController) GetStacks() error {

	stackResources, err := c.GetResourceList("/v2/stacks")
	resultMap := make(map[string]*StackInfo)

	for _, value := range *stackResources {
		si := new(StackInfo)
		si.Metadata = value.Metadata
		err = json.Unmarshal(value.Entity, &si.Entity)
		if err != nil {
			return err
		}
		resultMap[si.Entity.Name] = si
		// fmt.Println(si.Metadata.GUID)
		// fmt.Println(si.Metadata.CreatedAt)
		// fmt.Println(si.Metadata.UpdatedAt)
		// fmt.Println(si.Entity.Name)
		// fmt.Println(si.Entity.Description)
		// fmt.Println("============")
	}

	c.StackMap = &resultMap
	return err

}
