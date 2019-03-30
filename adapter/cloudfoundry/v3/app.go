package v3

// AppEntity
type App struct {
	GUID          string           `json:"guid"`       //"1cb006ee-fb05-47e1-b541-c34179ddc446"
	Name          string           `json:"name"`       //"my_app"
	State         string           `json:"state"`      //"STOPPED"
	CreatedAt     string           `json:"created_at"` //"2016-03-17T21:41:30Z"
	UpdatedAt     string           `json:"updated_at"` //"2016-06-08T16:41:26Z",
	Lifecycle     *LifecycleEntity `json:"lifecycle"`
	Relationships *Relationships   `json:"relationships"`
	Links         *Links           `json:"links"`
}

// LifecycleEntity
type LifecycleEntity struct {
	Type string         `json:"type"` //"buildpack"
	Data *LifecycleData `json:"data"`
}

// LifecycleData
type LifecycleData struct {
	Buildpacks []string `json:"buildpacks"` //["java_buildpack"]
	Stack      string   `json:"stack"`      //cflinuxfs2"
}

// Relationships
type Relationships struct {
	Space *RelationshipsSpace `json:"space"`
}

// RelationshipsSpace
type RelationshipsSpace struct {
	Data *SpaceData `json:"data"`
}

// SpaceData
type SpaceData struct {
	GUID string `json:"guid"` //"2f35885d-0c9d-4423-83ad-fd05066f8576"
}

// Links
type Links struct {
	Self                 *Link     `json:"self"`
	Space                *Link     `json:"space"`
	Processes            *Link     `json:"processes"`
	RouteMappings        *Link     `json:"route_mappings"`
	Packages             *Link     `json:"packages"`
	EnvironmentVariables *Link     `json:"environment_variables"`
	CurrentDroplet       *Link     `json:"current_droplet"`
	Droplets             *Link     `json:"droplets"`
	Tasks                *Link     `json:"tasks"`
	Start                *Link     `json:"start"`
	Stop                 *Link     `json:"stop"`
	Revisions            *Link     `json:"revisions"`
	DeployedRevisions    *Link     `json:"deployed_revisions"`
	Metadata             *Metadata `json:"metadata"`
}

// Link
type Link struct {
	HRef   string `json:"href"`
	Method string `json:"method"`
}

// Metadata
type Metadata struct {
	Labels      string `json:"labels"`
	Annotations string `json:"annotations"`
}
