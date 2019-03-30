package cloudfoundry

import (
	"context"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCloudControllerCreation(t *testing.T) {

	Convey("Given an config without username", t, func() {

		Convey("When CloudController is created", func() {

			Convey("Then an error message indicates the missing username", func() {
				cloudcontrollerConfig := CloudControllerConfig{
					Username: "",
				}

				_, err := NewCloudController(cloudcontrollerConfig)

				So(err.Error(), ShouldEqual, "username cannot be empty")
			})

		})

	})

	Convey("Given a config without password", t, func() {

		Convey("When CloudController is created", func() {

			Convey("Then an error message indicates the missing password", func() {
				cloudcontrollerConfig := CloudControllerConfig{
					Username: "cloudmaster",
					Password: "",
				}

				_, err := NewCloudController(cloudcontrollerConfig)

				So(err.Error(), ShouldEqual, "password cannot be empty")
			})

		})

	})

	Convey("Given a config without apiurl", t, func() {

		Convey("When CloudController is created", func() {

			Convey("Then an error message indicates the missing apiurl", func() {
				cloudcontrollerConfig := CloudControllerConfig{
					Username:     "cloudmaster",
					Password:     "cloudpass",
					APIURLString: "",
				}

				_, err := NewCloudController(cloudcontrollerConfig)

				So(err.Error(), ShouldEqual, "apiUrl cannot be empty")
			})

		})

	})

	Convey("Given a config with malformed apiurl", t, func() {

		Convey("When CloudController is created", func() {

			Convey("Then an error message indicates the malformed apiurl", func() {
				cloudcontrollerConfig := CloudControllerConfig{
					Username:     "cloudmaster",
					Password:     "cloudpass",
					APIURLString: ":8080/hhh//",
				}

				_, err := NewCloudController(cloudcontrollerConfig)

				So(err.Error(), ShouldEqual, "parse :8080/hhh//: missing protocol scheme")
			})

		})

	})

	Convey("Given a config with username, password and apiurl", t, func() {

		Convey("When CloudController is created", func() {

			Convey("Then a new instance of the CloudController is returned without error", func() {
				cloudcontrollerConfig := CloudControllerConfig{
					Username:     "cloudmaster",
					Password:     "cloudpass",
					APIURLString: "http://localhost",
				}

				cc, err := NewCloudController(cloudcontrollerConfig)

				So(err, ShouldEqual, nil)
				So(cc, ShouldNotEqual, nil)
			})

		})

	})

	Convey("Given a config with username, password and apiurl", t, func() {

		Convey("When CloudController Info is requested", func() {

			Convey("Then a valid version 2 info object is returned", func(c C) {
				cloudcontrollerConfig := CloudControllerConfig{
					Username:     "cloudmaster",
					Password:     "cloudpass",
					APIURLString: "http://api.mycloudcontroller",
				}

				okResponse := `{
					"name": "",
					"build": "",
					"support": "",
					"version": 0,
					"description": " ",
					"authorization_endpoint": " ",
					"token_endpoint": " ",
					"min_cli_version": "6.22.0",
					"min_recommended_cli_version": "latest",
					"app_ssh_endpoint": " ",
					"app_ssh_host_key_fingerprint": " ",
					"app_ssh_oauth_client": " ",
					"doppler_logging_endpoint": " ",
					"api_version": "2.133.0",
					"osbapi_version": "2.14",
					"routing_endpoint": " "
				}`

				h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					c.So(r.Header.Get("accept"), ShouldEqual, "application/json")
					c.So(r.Method, ShouldEqual, "GET")
					w.Write([]byte(okResponse))
				})
				httpClient, teardown := testingHTTPClient(h)
				defer teardown()

				cc, _ := NewCloudController(cloudcontrollerConfig)
				cc.httpClient = httpClient

				v2info, err := cc.GetV2Info()
				So(err, ShouldEqual, nil)
				So(v2info, ShouldNotEqual, nil)
				So(v2info.APIVersion, ShouldEqual, "2.133.0")
			})

		})

	})

	Convey("Given a config with username, password and apiurl", t, func() {

		Convey("When CloudController Login is requested", func() {

			Convey("Then a valid oauth token is returned", func(c C) {
				cloudcontrollerConfig := CloudControllerConfig{
					Username:     "cloudmaster",
					Password:     "cloudpass",
					APIURLString: "http://api.mycloudcontroller",
				}

				okInfoResponse := `{
					"name": "",
					"build": "",
					"support": "",
					"version": 0,
					"description": " ",
					"authorization_endpoint": "http://api.mycloudcontroller",
					"token_endpoint": "http://api.mycloudcontroller",
					"min_cli_version": "6.22.0",
					"min_recommended_cli_version": "latest",
					"app_ssh_endpoint": " ",
					"app_ssh_host_key_fingerprint": " ",
					"app_ssh_oauth_client": " ",
					"doppler_logging_endpoint": " ",
					"api_version": "2.133.0",
					"osbapi_version": "2.14",
					"routing_endpoint": " "
				}`

				okResponse := `{
					"access_token": "eyhuetzligruetzli",
					"token_type": "bearer",
					"id_token": "eyhuetzligruetzli.eyhuetzligruetzli.eyhuetzligruetzliid",
					"refresh_token": "eyhuetzligruetzli.eyhuetzligruetzli.eyhuetzligruetzlirefresh",
					"expires_in": 300,
					"scope": "openid uaa.user cloud_controller.read password.write cloud_controller.write",
					"jti": "1234567890"
				}`

				h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

					if r.Method == "GET" {
						w.Write([]byte(okInfoResponse))
					}

					if r.Method == "POST" {
						c.So(r.Header.Get("authorization"), ShouldEqual, "Basic Y2Y6")
						c.So(r.Header.Get("content-type"), ShouldEqual, "application/x-www-form-urlencoded")
						b, _ := ioutil.ReadAll(r.Body)
						bodystring := string(b[:])
						c.So(bodystring, ShouldEqual, "grant_type=password&password=cloudpass&scope=&username=cloudmaster")
						w.Write([]byte(okResponse))
					}

				})
				httpClient, teardown := testingHTTPClient(h)
				defer teardown()

				cc, _ := NewCloudController(cloudcontrollerConfig)
				cc.httpClient = httpClient

				err := cc.Login()
				So(err, ShouldEqual, nil)
				So(cc.AccessToken.AccessToken, ShouldEqual, "eyhuetzligruetzli")
			})

		})

	})

}

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return cli, s.Close
}
