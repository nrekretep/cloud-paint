package services

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSingleAppDiagram(t *testing.T) {

	Convey("Given the config for the diagram does not exist", t, func() {

		Convey("When the SingleAppDiagram Service is created", func() {

			Convey("Then an error messages indicates the missing config", func() {
				singleAppDiagramService, err := NewSingleAppDiagramService(nil)

				So(singleAppDiagramService, ShouldEqual, nil)
				So(err.Error(), ShouldEqual, "a non empty config must be provided to a diagram service")
			})

		})

	})

	Convey("Given the config for the diagram exists", t, func() {

		Convey("When the SingleAppDiagram Service is created", func() {

			Convey("Then a valid diagram service is returned", func() {
				config := Config{Usename: "u", Password: "p", ApiUrl: "test"}
				singleAppDiagramService, err := NewSingleAppDiagramService(&config)

				So(singleAppDiagramService, ShouldNotEqual, nil)
				So(err, ShouldEqual, nil)
			})

		})

	})

	Convey("Given the ID for the app is empty", t, func() {

		Convey("When the SingleAppDiagram Service is created", func() {

			Convey("Then an error messages indicates the empty app ID", func() {
				config := Config{Usename: "u", Password: "p", ApiUrl: "test"}
				singleAppDiagramService, err := NewSingleAppDiagramService(&config)

				So(singleAppDiagramService, ShouldNotEqual, nil)
				So(err, ShouldEqual, nil)

				appID := ""
				_, err = singleAppDiagramService.GetRawDiagram(appID)
				So(err, ShouldNotEqual, nil)
				So(err.Error(), ShouldEqual, "a valid id for the app must be provided")
			})

		})

	})

	Convey("Given app with the given ID does not exist", t, func() {

		Convey("When the SingleAppDiagram is rendered", func() {

			Convey("Then an error messages indicates the wrong app ID", nil)

		})

	})

}
