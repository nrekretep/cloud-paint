package services

import (
	"errors"
	//"fmt"
	"github.com/nrekretep/cloudpaint/adapter/cloudfoundry"
	"github.com/nrekretep/cloudpaint/adapter/plantuml"
)

// SingleAppDiagramService -
type SingleAppDiagramService struct {
	config *Config
}

// NewSingleAppDiagramService -
func NewSingleAppDiagramService(c *Config) (*SingleAppDiagramService, error) {

	if c == nil {
		return nil, errors.New("a non empty config must be provided to a diagram service")
	}

	diagramService := &SingleAppDiagramService{config: c}

	return diagramService, nil
}

// renderTemplate -
func (s *SingleAppDiagramService) GetRawDiagram(appID string) (string, error) {

	if appID == "" {
		return "", errors.New("a valid id for the app must be provided")
	}

	cloudControllerConfig := cloudfoundry.CloudControllerConfig{Username: s.config.Usename, Password: s.config.Password, APIURLString: s.config.ApiUrl}
	cloudController, err := cloudfoundry.NewCloudController(cloudControllerConfig)

	if err != nil {
		return "", err
	}

	err = cloudController.Login()
	if err != nil {
		return "", err
	}

	app, err := cloudController.GetV3App(appID)

	if err != nil {
		return "", err
	}

	cloudController.GetOrganizations()
	if err != nil {
		return "", err
	}

	cloudController.GetSpaces()
	if err != nil {
		return "", err
	}

	cloudController.GetStacks()
	if err != nil {
		return "", err
	}

	cloudController.GetBuildpacks()
	if err != nil {
		return "", err
	}

	plantUml := plantuml.NewPlantUML(cloudController)

	return plantUml.CreateSingleAppDiagram(app), nil

}
