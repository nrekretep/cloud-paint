package plantuml

import (
	"github.com/nrekretep/cloudpaint/adapter/cloudfoundry"
	"github.com/nrekretep/cloudpaint/adapter/cloudfoundry/v3"
	"strings"
)

type PlantUML struct {
	CloudController *cloudfoundry.CloudController
}

// CreateDiagram -
func (p *PlantUML) CreateDiagram() string {
	var stringBuilder strings.Builder

	p.WriteStartTag(&stringBuilder)

	p.WriteAllStacks(&stringBuilder)
	p.WriteAllBuildpacks(&stringBuilder)
	p.WriteBuildpackStackRelation(&stringBuilder)

	p.WriteAllOrgs(&stringBuilder)
	p.WriteAllSpaces(&stringBuilder)
	p.WriteAllOrgSpaceRelations(&stringBuilder)

	p.WriteAllApps(&stringBuilder)
	p.WriteSpaceAppRelation(&stringBuilder)
	p.WriteAllAppBuildpackRelation(&stringBuilder)
	p.WriteEndTag(&stringBuilder)

	return stringBuilder.String()
}

// CreateSingleAppDiagram -
func (p *PlantUML) CreateSingleAppDiagram(app *v3.App) string {
	var stringBuilder strings.Builder

	p.WriteStartTag(&stringBuilder)
	p.WriteSkin(&stringBuilder)

	p.WriteTitle(&stringBuilder, "Single App Diagram - "+app.Name)

	space := (*p.CloudController.SpaceMap)[app.Relationships.Space.Data.GUID]
	p.WriteSpace(&stringBuilder, space)

	org := (*p.CloudController.OrganizationMap)[space.Entity.OrganizationGUID]
	p.WriteOrg(&stringBuilder, org)

	p.WriteOrgSpaceRelation(&stringBuilder, org.Metadata.GUID, space.Metadata.GUID)

	p.WriteApp(&stringBuilder, app)

	if app.Lifecycle.Type == "buildpack" {

		for _, b := range app.Lifecycle.Data.Buildpacks {
			p.WriteBuildpack(&stringBuilder, b)
			p.WriteAppBuildpackRelation(&stringBuilder, app, b)
		}

		p.WriteStack(&stringBuilder, app.Lifecycle.Data.Stack)
		p.WriteAppStackRelation(&stringBuilder, app, app.Lifecycle.Data.Stack)
	}

	p.WriteAppSpaceRelation(&stringBuilder, app)

	p.WriteEndTag(&stringBuilder)

	return stringBuilder.String()
}

// WriteAppStackRelation -
func (p *PlantUML) WriteAppStackRelation(sb *strings.Builder, app *v3.App, s string) {

	sb.WriteString(*p.TrimGUID(&app.GUID))
	sb.WriteString(" --> ")
	sb.WriteString(s)
	sb.WriteString("\n")

}

// WriteStack -
func (p *PlantUML) WriteStack(sb *strings.Builder, s string) {

	sb.WriteString("[**")
	sb.WriteString(s)
	sb.WriteString("**] <<stack>> as ")
	sb.WriteString(s)
	sb.WriteString("\n")

}

// WriteAppBuildpackRelation -
func (p *PlantUML) WriteAppBuildpackRelation(sb *strings.Builder, app *v3.App, b string) {

	sb.WriteString(*p.TrimGUID(&app.GUID))
	sb.WriteString(" --> ")
	sb.WriteString(b)
	sb.WriteString("\n")

}

// WriteBuildpack -
func (p *PlantUML) WriteBuildpack(sb *strings.Builder, b string) {

	sb.WriteString("[**")
	sb.WriteString(b)
	sb.WriteString("**] <<buildpack>> as ")
	sb.WriteString(b)
	sb.WriteString("\n")

}

// WriteApp -
func (p *PlantUML) WriteApp(sb *strings.Builder, app *v3.App) {

	sb.WriteString("component ")
	sb.WriteString(*p.TrimGUID(&app.GUID))
	sb.WriteString(" <<app>> [\n**")
	sb.WriteString(app.Name)
	sb.WriteString("**\n")
	sb.WriteString("State: " + app.State + "\n")
	sb.WriteString("Created at: " + app.CreatedAt + "\n")
	sb.WriteString("Updated at: " + app.UpdatedAt + "\n")
	sb.WriteString("]")
	sb.WriteString("\n")

}

// WriteAppSpaceRelation -
func (p *PlantUML) WriteAppSpaceRelation(sb *strings.Builder, app *v3.App) {

	sb.WriteString(*p.TrimGUID(&app.Relationships.Space.Data.GUID))
	sb.WriteString(" --> ")
	sb.WriteString(*p.TrimGUID(&app.GUID))
	sb.WriteString("\n")

}

// WriteOrgSpaceRelation -
func (p *PlantUML) WriteOrgSpaceRelation(sb *strings.Builder, orgGUID string, spaceGUID string) {

	sb.WriteString(*p.TrimGUID(&orgGUID))
	sb.WriteString(" --> ")
	sb.WriteString(*p.TrimGUID(&spaceGUID))
	sb.WriteString("\n")

}

// WriteOrg -
func (p *PlantUML) WriteOrg(sb *strings.Builder, org *cloudfoundry.OrganizationInfo) {

	sb.WriteString("[**")
	sb.WriteString(org.Entity.Name)
	sb.WriteString("**] <<organization>> as ")
	sb.WriteString(*p.TrimGUID(&org.Metadata.GUID))
	sb.WriteString("\n")

}

// WriteSpace -
func (p *PlantUML) WriteSpace(sb *strings.Builder, space *cloudfoundry.SpaceInfo) {

	sb.WriteString("[**")
	sb.WriteString(space.Entity.Name)
	sb.WriteString("**] <<space>> as ")
	sb.WriteString(*p.TrimGUID(&space.Metadata.GUID))
	sb.WriteString("\n")

}

// NewPlantUML -
func NewPlantUML(c *cloudfoundry.CloudController) *PlantUML {

	plantUML := &PlantUML{CloudController: c}

	return plantUML
}

// WriteAllStacks -
func (p *PlantUML) WriteAllStacks(sb *strings.Builder) {

	for _, v := range *p.CloudController.StackMap {

		sb.WriteString("[")
		sb.WriteString(v.Entity.Name)
		sb.WriteString("] <<stack>> as ")
		sb.WriteString(v.Entity.Name)
		sb.WriteString("\n")

	}
	sb.WriteString("\n")
}

// WriteAllBuildpacks -
func (p *PlantUML) WriteAllBuildpacks(sb *strings.Builder) {

	for _, v := range *p.CloudController.BuildpackMap {

		sb.WriteString("[")
		sb.WriteString(v.Entity.Name)
		sb.WriteString("] <<buildpack>> as ")
		sb.WriteString(*p.TrimGUID(&v.Metadata.GUID))
		sb.WriteString("\n")

	}
	sb.WriteString("\n")
}

// WriteBuildpackStackRelation -
func (p *PlantUML) WriteBuildpackStackRelation(sb *strings.Builder) {

	for _, v := range *p.CloudController.BuildpackMap {

		if v.Entity.Stack != "" {
			sb.WriteString(*p.TrimGUID(&v.Metadata.GUID))
			sb.WriteString(" --> ")
			sb.WriteString((*p.CloudController.StackMap)[v.Entity.Stack].Entity.Name)
			sb.WriteString("\n")
		}
	}
	sb.WriteString("\n")
}

// WriteAllOrgs -
func (p *PlantUML) WriteAllOrgs(sb *strings.Builder) {

	for _, v := range *p.CloudController.OrganizationMap {

		sb.WriteString("[")
		sb.WriteString(v.Entity.Name)
		sb.WriteString("] <<org>> as ")
		sb.WriteString(*p.TrimGUID(&v.Metadata.GUID))
		sb.WriteString("\n")

	}
	sb.WriteString("\n")
}

// WriteAllSpaces -
func (p *PlantUML) WriteAllSpaces(sb *strings.Builder) {

	for _, v := range *p.CloudController.SpaceMap {

		sb.WriteString("[")
		sb.WriteString(v.Entity.Name)
		sb.WriteString("] <<space>> as ")
		sb.WriteString(*p.TrimGUID(&v.Metadata.GUID))
		sb.WriteString("\n")

	}
	sb.WriteString("\n")
}

// WriteAllOrgSpaceRelations -
func (p *PlantUML) WriteAllOrgSpaceRelations(sb *strings.Builder) {

	for _, v := range *p.CloudController.SpaceMap {

		sb.WriteString(*p.TrimGUID(&(*p.CloudController.OrganizationMap)[v.Entity.OrganizationGUID].Metadata.GUID))
		sb.WriteString(" --> ")
		sb.WriteString(*p.TrimGUID(&v.Metadata.GUID))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
}

// WriteAllApps -
func (p *PlantUML) WriteAllApps(sb *strings.Builder) {

	for _, v := range *p.CloudController.AppMap {

		sb.WriteString("[")
		sb.WriteString(v.Entity.Name)
		sb.WriteString("] <<app>> as ")
		sb.WriteString(*p.TrimGUID(&v.Metadata.GUID))
		sb.WriteString("\n")

	}
	sb.WriteString("\n")
}

// WriteSpaceAppRelation -
func (p *PlantUML) WriteSpaceAppRelation(sb *strings.Builder) {

	for _, v := range *p.CloudController.AppMap {

		sb.WriteString(*p.TrimGUID(&(*p.CloudController.SpaceMap)[v.Entity.SpaceGUID].Metadata.GUID))
		sb.WriteString(" --> ")
		sb.WriteString(*p.TrimGUID(&v.Metadata.GUID))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
}

// WriteAppBuildpackRelation -
func (p *PlantUML) WriteAllAppBuildpackRelation(sb *strings.Builder) {

	for _, v := range *p.CloudController.AppMap {

		if v.Entity.DetectedBuildpackGUID != "" {
			sb.WriteString(*p.TrimGUID(&v.Metadata.GUID))
			sb.WriteString(" --> ")
			sb.WriteString(*p.TrimGUID(&(*p.CloudController.BuildpackMap)[v.Entity.DetectedBuildpackGUID].Metadata.GUID))
		}

		sb.WriteString("\n")
	}
	sb.WriteString("\n")
}

func (p *PlantUML) TrimGUID(guid *string) *string {
	t := strings.Replace(*guid, "-", "", -1)
	return &t
}

// WriteStartTag -
func (p *PlantUML) WriteStartTag(sb *strings.Builder) {
	sb.WriteString("@startuml\n")
}

// WriteEndTag -
func (p *PlantUML) WriteEndTag(sb *strings.Builder) {
	sb.WriteString("center footer Generated with cloudpaint (https://github.com/nrekretep/cloudpaint)\n")
	sb.WriteString("@enduml\n")
}

// WriteTitle -
func (p *PlantUML) WriteTitle(sb *strings.Builder, diagramtitle string) {
	sb.WriteString("title " + diagramtitle + "\n")
}

// WriteSkin -
func (p *PlantUML) WriteSkin(sb *strings.Builder) {
	sb.WriteString(`skinparam componentStyle uml2

skinparam component {
	
FontSize 18
FontName Impact
FontColor #009F9D
	
StereotypeFontName Impact
StereotypeFontSize 14
StereotypeFontColor #0f0a3c
	
BorderColor #0F0A3C
	
BackgroundColor #cdffeb
	
ArrowFontName Impact
ArrowColor #0F0A3C
ArrowFontColor #777777
}
	
skinparam titleBorderRoundCorner 5
skinparam titleBorderThickness 2
skinparam titleBorderColor #393e46
skinparam titleBackgroundColor #eeeeee
	
skinparam footerFontColor #07456f`)
	sb.WriteString("\n")
}
