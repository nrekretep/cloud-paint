package plantuml

import (
	"github.com/nrekretep/cloudpaint/adapter/cloudfoundry"
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
	p.WriteOrgSpaceRelation(&stringBuilder)

	p.WriteAllApps(&stringBuilder)
	p.WriteSpaceAppRelation(&stringBuilder)
	p.WriteAppBuildpackRelation(&stringBuilder)
	p.WriteEndTag(&stringBuilder)

	return stringBuilder.String()
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

// WriteOrgSpaceRelation -
func (p *PlantUML) WriteOrgSpaceRelation(sb *strings.Builder) {

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
func (p *PlantUML) WriteAppBuildpackRelation(sb *strings.Builder) {

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
	sb.WriteString("@enduml\n")
}
