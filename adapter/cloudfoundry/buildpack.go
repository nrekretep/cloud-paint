package cloudfoundry

import (
	"encoding/json"
	"fmt"
)

// BuildpackEntity - Entity data for buildpacks.
type BuildpackEntity struct {
	Name     string `json:"name"`
	Stack    string `json:"stack"`
	Position int    `json:"position"`
	Enabled  bool   `json:"enabled"`
	Locked   bool   `json:"locked"`
	Filename string `json:"filename"`
}

// BuildpackInfo - Details about Buildpacks.
type BuildpackInfo struct {
	Metadata Metadata
	Entity   BuildpackEntity
}

// GetBuildpacks - Loads infos about all buildpacks
func (c *CloudController) GetBuildpacks() error {

	buildpackResources, err := c.GetResourceList("/v2/buildpacks")

	resultMap := make(map[string]*BuildpackInfo)

	for _, value := range *buildpackResources {
		bpi := new(BuildpackInfo)
		bpi.Metadata = value.Metadata
		err = json.Unmarshal(value.Entity, &bpi.Entity)

		if err != nil {
			fmt.Println(err)
			return err
		}

		resultMap[bpi.Metadata.GUID] = bpi
		// fmt.Println(bpi.Metadata.GUID)
		// fmt.Println(bpi.Metadata.CreatedAt)
		// fmt.Println(bpi.Metadata.UpdatedAt)
		// fmt.Println(bpi.Entity.Name)
		// fmt.Println(bpi.Entity.Filename)
		// fmt.Println(bpi.Entity.Position)
		// fmt.Println(bpi.Entity.Locked)
		// fmt.Println(bpi.Entity.Enabled)
		// fmt.Println("============")
	}

	c.BuildpackMap = &resultMap
	return err

}
