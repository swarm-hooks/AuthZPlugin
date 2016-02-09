package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"net/http/httptest"

	log "github.com/Sirupsen/logrus"

	"regexp"
	"strconv"

	"github.com/AuthZPluginBackEnd/headers"
	"github.com/gorilla/mux"
	"github.com/jeffail/gabs"
	//	"encoding/json"
	"github.com/samalba/dockerclient"
)

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

///////////////////////////////////DEprecate//////////////////////////

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

//UTILS

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

/*
func CheckLinksOwnerShip(cluster cluster.Cluster, tenantName string, containerConfig dockerclient.ContainerConfig) (bool, *ValidationOutPutDTO) {
	log.Debug("in CheckLinksOwnerShip")
	log.Debugf("containerConfig: %+v", containerConfig)
	linksSize := len(containerConfig.HostConfig.Links)
	volumesFrom := containerConfig.HostConfig.VolumesFrom
	volumesFromSize := len(containerConfig.HostConfig.VolumesFrom)
	if linksSize < 1 && volumesFromSize < 1 {
		return true, &ValidationOutPutDTO{ContainerID: ""}

	}
	log.Debug("Checking links...")
	containers := cluster.Containers()
	linkSet := make(map[string]bool)
	links := make([]string, 0)
	var v int // count of volumesFrom links validated
	var l int // count of links validated

	log.Debug("**************************************************")
	for _, container := range containers {
		if strings.HasSuffix(container.Info.Name, tenantName) {
			log.Debugf("Examine container %s %s", container.Info.Name, container.Info.Id)
			for i := 0; i < volumesFromSize; i++ {
				if v == volumesFromSize {
					break
				}
				log.Debugf("Examine VolumeFrom[%d] == %s", i, containerConfig.HostConfig.VolumesFrom[i])
				// volumesFrom element format <container_name>:<RW|RO>
				volumeFromArray := strings.SplitN(strings.TrimSpace(containerConfig.HostConfig.VolumesFrom[i]), ":", 2)
				volumeFrom := strings.TrimSpace(volumeFromArray[0])
				if strings.HasPrefix(container.Info.Id, volumeFrom) {
					log.Debug("volumesFrom element with container id matches tenant container")
					// no need to modify volumesFrom
					v++
				} else if container.Info.Name == "/"+volumeFrom+tenantName {
					log.Debug("volumesFrom element with container name matches tenant container")
					volumesFrom[i] = container.Info.Id
					if len(volumeFromArray) > 1 {
						volumesFrom[i] += ":"
						volumesFrom[i] += strings.TrimSpace(volumeFromArray[1])
					}
					v++
				}

			}
			for i := 0; i < linksSize; i++ {
				if l == linksSize {
					break
				}
				log.Debugf("Examine links[%d] == %s", i, containerConfig.HostConfig.Links[i])

				linkArray := strings.SplitN(containerConfig.HostConfig.Links[i], ":", 2)
				link := strings.TrimSpace(linkArray[0])
				if strings.HasPrefix(container.Info.Id, link) || "/"+link+tenantName == container.Info.Name {
					log.Debug("Add link and alias to linkset")
					_, ok := linkSet[link]
					if !ok {
						linkSet[link] = true
						links = append(links, container.Info.Id+":"+link)
					}
					// check for alias
					if len(linkArray) > 1 {
						links = append(links, container.Info.Id+":"+strings.TrimSpace(linkArray[1]))
					}
					l++
				}
			}
		}
		// are we done?
		if v == volumesFromSize && l == linksSize {
			break
		}
	}

	if v != volumesFromSize || l != linksSize {
		return false, &ValidationOutPutDTO{ContainerID: "", ErrorMessage: "Tenant does not own containers in volumesFrom or links."}
	}
	return true, &ValidationOutPutDTO{ContainerID: "", Links: links, VolumesFrom: volumesFrom}

}
*/

type Config struct {
	HostConfig struct {
		Links       []interface{}
		VolumesFrom []interface{}
	}
}

/*
//TODO - Pass by ref ?
func CheckOwnerShip(cluster cluster.Cluster, tenantName string, r *http.Request) (bool, *ValidationOutPutDTO) {
	containers := cluster.Containers()
	log.Debug("got name: ", mux.Vars(r)["name"])
	if mux.Vars(r)["name"] == "" {
		return true, nil
	}
	tenantSet := make(map[string]bool)
	for _, container := range containers {
		if "/"+mux.Vars(r)["name"]+tenantName == container.Info.Name {
			log.Debug("Match By name!")
			return true, &ValidationOutPutDTO{ContainerID: container.Info.Id, Links: nil}
		} else if "/"+mux.Vars(r)["name"] == container.Info.Name {
			if container.Labels[headers.TenancyLabel] == tenantName {
				return true, &ValidationOutPutDTO{ContainerID: container.Info.Id, Links: nil}
			}
		} else if mux.Vars(r)["name"] == container.Info.Id {
			log.Debug("Match By full ID! Checking Ownership...")
			log.Debug("Tenant name: ", tenantName)
			log.Debug("Tenant Lable: ", container.Labels[headers.TenancyLabel])
			if container.Labels[headers.TenancyLabel] == tenantName {
				return true, &ValidationOutPutDTO{ContainerID: container.Info.Id, Links: nil}
			}
			return false, nil

		}
		if container.Labels[headers.TenancyLabel] == tenantName {
			tenantSet[container.Id] = true
		}
	}

	//Handle short ID
	ambiguityCounter := 0
	var returnID string
	for k := range tenantSet {
		if strings.HasPrefix(cluster.Container(k).Info.Id, mux.Vars(r)["name"]) {
			ambiguityCounter++
			returnID = cluster.Container(k).Info.Id
		}
		if ambiguityCounter == 1 {
			log.Debug("Matched by short ID")
			return true, &ValidationOutPutDTO{ContainerID: returnID, Links: nil}
		}
		if ambiguityCounter > 1 {
			log.Debug("Ambiguiy by short ID")
			//TODO - ambiguity
		}
		if ambiguityCounter == 0 {
			log.Debug("No match by short ID")
			//TODO - no such container
		}
	}
	return false, nil
}
*/
func CleanUpLabeling(r *http.Request, rec *httptest.ResponseRecorder) []byte {
	newBody := bytes.Replace(rec.Body.Bytes(), []byte(headers.TenancyLabel), []byte(" "), -1)
	//TODO - Here we just use the token for the tenant name for now so we remove it from the data before returning to user.
	newBody = bytes.Replace(newBody, []byte(r.Header.Get(headers.AuthZTenantIdHeaderName)), []byte(" "), -1)
	newBody = bytes.Replace(newBody, []byte(",\" \":\" \""), []byte(""), -1)
	log.Debug("Got this new body...", string(newBody))
	return newBody
}

func ParseField(field string, fieldType interface{}, body []byte) (interface{}, error) {
	log.Debugf("In parseField, field: %s Request body: %s", field, string(body))
	jsonParsed, err := gabs.ParseJSON(body)
	if err != nil {
		log.Error("failed to parse!")
		return nil, err
	}

	switch v := fieldType.(type) {
	case float64:
		log.Debug("Parsing type: ", v)
		parsedField, ok := jsonParsed.Path(field).Data().(float64)
		if ok {
			res := strconv.FormatFloat(parsedField, 'f', -1, 64)
			log.Debugf("Parsed field: " + res)
			return parsedField, nil
		}
	case []string:
		log.Debug("Parsing type: ", v)
		parsedField, ok := jsonParsed.Path(field).Data().([]string)
		if ok {
			log.Debug(parsedField)
			return parsedField, nil
		}
	default:
		log.Error("Unknown field type to parse")
	}

	return nil, errors.New(fmt.Sprintf("failed to parse field %s from request body %s", field, string(body)))
}
