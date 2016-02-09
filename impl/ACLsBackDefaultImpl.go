package impl

import (
	//	"bytes"
	//	"io/ioutil"

	"strings"

	//	"github.com/docker/swarm/cluster/swarm"
	"github.com/AuthZPluginBackEnd/states"
	log "github.com/Sirupsen/logrus"
	//	"github.com/docker/swarm/cluster/swarm"
	"github.com/AuthZPluginBackEnd/authz"
	"github.com/AuthZPluginBackEnd/headers"
	"github.com/AuthZPluginBackEnd/utils"
	"github.com/samalba/dockerclient"
)

//DefaultACLsImpl - Default implementation of ACLs API
type ACLsBackDefaultImpl struct{}

/*
ValidateRequest - Who wants to do what - allow or not
*/

func (*ACLsBackDefaultImpl) ValidateRequest(req authz.Request, containerConfig dockerclient.ContainerConfig) bool {
	eventType := eventParse(req.RequestURI)

	switch eventType {
	case states.ContainerCreate:
		tenantIdToValidate := req.RequestHeaders[headers.AuthZTenantIdHeaderName]
		//		tenantIdToValidate := r.Header.Get(headers.AuthZTenantIdHeaderName)
		isValid := utils.CheckRefsOwnerships(tenantIdToValidate, containerConfig)
		return isValid

	case states.ContainersList:
		//Take action on the response...
		return true

	case states.Unauthorized:
		return false

	default:
		//CONTAINER_INSPECT / CONTAINER_OTHERS / STREAM_OR_HIJACK / PASS_AS_IS
		tenantIdToValidate := req.RequestHeaders[headers.AuthZTenantIdHeaderName]
		return utils.CheckOwnerShip(tenantIdToValidate, req.RequestURI, containerConfig)

	}
}

//startRegExp = regexp.MustCompile(`/containers/(.*)/start$`)
//createRegExp = regexp.MustCompile(`/containers/(.*)/create$`)
//	listRegExp = regexp.MustCompile(`/containers/(.*)/json(.*)`)

/*Probably should use regular expressions here*/
func eventParse(reqURI string) states.EventEnum {
	log.Debug("Got the uri...", reqURI)

	if strings.Contains(reqURI, "/containers") && (strings.Contains(reqURI, "create")) {
		return states.ContainerCreate
	}

	if strings.Contains(reqURI, "/containers/json") {
		return states.ContainersList
	}

	if strings.Contains(reqURI, "/containers") &&
		(strings.Contains(reqURI, "logs") || strings.Contains(reqURI, "attach") || strings.HasSuffix(reqURI, "exec")) {
		return states.StreamOrHijack
	}
	if strings.Contains(reqURI, "/containers") && strings.HasSuffix(reqURI, "/json") {
		return states.ContainerInspect
	}
	if strings.Contains(reqURI, "/containers") {
		return states.ContainerOthers
	}
	if strings.Contains(reqURI, "/images") && strings.HasSuffix(reqURI, "/json") {
		return states.PassAsIs
	}
	if strings.HasSuffix(reqURI, "/version") || strings.Contains(reqURI, "/exec/") {
		return states.PassAsIs
	}
	return states.NotSupported
}

//Init - Any required initialization
//func (*DefaultACLsImpl) Init() error {
//	return nil
//}
