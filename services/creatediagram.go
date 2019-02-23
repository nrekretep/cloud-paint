package services

import (
	"github.com/nrekretep/cloudpaint/adapter/cloudfoundry"
	"github.com/nrekretep/cloudpaint/adapter/plantuml"
)

// CreateDiagramService -
type CreateDiagramService struct {
	CloudController *cloudfoundry.CloudController
}

// NewCreateDiagramService -
func NewCreateDiagramService(c *cloudfoundry.CloudController) *CreateDiagramService {

	diagramService := &CreateDiagramService{CloudController: c}

	return diagramService
}

// renderTemplate -
func (c *CreateDiagramService) RenderTemplate() string {

	plantUml := plantuml.NewPlantUML(c.CloudController)

	return plantUml.CreateDiagram()
}
