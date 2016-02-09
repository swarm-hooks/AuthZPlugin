package api

import (
	"github.com/AuthZPluginBackEnd/authz"
	"github.com/samalba/dockerclient"
)

type ACLsBackAPI interface {
	ValidateRequest(req authz.Request, containerConfig dockerclient.ContainerConfig) bool
}
