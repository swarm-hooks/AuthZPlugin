package utils

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"regexp"

	"github.com/gorilla/mux"
	//	"encoding/json"
	"github.com/samalba/dockerclient"
)

//UTILS

//For create check links etc
func CheckOwnerShip(tenantIdToValidate string, reqURI string, containerConfig dockerclient.ContainerConfig) bool {
	re := regexp.MustCompile(`/containers/(.*)/(.*)`)
	idCmdArr := re.FindStringSubmatch(reqURI)
	if tenantIdToValidate != getTenantByIdOrName(idCmdArr[1]) {
		return false
	}

	return CheckRefsOwnerships(tenantIdToValidate, containerConfig)

}

func CheckRefsOwnerships(tenantIdToValidate string, containerConfig dockerclient.ContainerConfig) bool {

	return checkLinks(tenantIdToValidate, containerConfig.HostConfig.Links) &&
		checkVolumes(tenantIdToValidate, containerConfig.HostConfig.VolumesFrom)
}

func getTenantByIdOrName(idOrName string) string {
	//	dataStore	...
	//Key can be name, ID or short ID...
	return "space1"
}

func checkLinks(tenantIdToValidate string, links []string) bool {
	return false
}

func checkVolumes(tenantIdToValidate string, volumes []string) bool {
	return false
}

func ModifyRequest(r *http.Request, body io.Reader, urlStr string, containerID string) (*http.Request, error) {

	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
		r.Body = rc
	}
	if urlStr != "" {
		u, err := url.Parse(urlStr)

		if err != nil {
			return nil, err
		}
		r.URL = u
		mux.Vars(r)["name"] = containerID
	}

	return r, nil
}

type Config struct {
	HostConfig struct {
		Links       []interface{}
		VolumesFrom []interface{}
	}
}

//////////////////////////////////DEprecate//////////////////////////

type ValidationOutPutDTO struct {
	ContainerID string
	//	Links       map[string]string
	//	Links       map[string][]string
	Links []string
	//	VolumesFrom map[string]string
	VolumesFrom  []string
	ErrorMessage string
	//Quota can live here too? Currently quota needs only raise error
	//What else
}
