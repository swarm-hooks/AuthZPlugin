package main

import (
	"bytes"
	"encoding/json"
	"regexp"

	//	dockerapi "github.com/docker/docker/api"
	//	dockerclient "github.com/docker/engine-api/client"
	"github.com/AuthZPlugin/api"
	"github.com/AuthZPlugin/authz"
	"github.com/AuthZPlugin/impl"
	dockerclient "github.com/samalba/dockerclient"

	"github.com/docker/go-connections/tlsconfig"
)

var startRegExp = regexp.MustCompile(`/containers/(.*)/start$`)
var aclsAPI api.ACLsBackAPI

type authzPlugin struct {
	client *dockerclient.DockerClient
}

func newPlugin(dockerHost string) (*authzPlugin, error) {
	c, _ := tlsconfig.Client(tlsconfig.Options{InsecureSkipVerify: true})
	client, err := dockerclient.NewDockerClient(dockerHost, c)

	if err != nil {
		return nil, err
	}

	aclsAPI = new(impl.ACLsBackDefaultImpl)
	return &authzPlugin{client: client}, nil
}

//Before passing request to manager
func (p *authzPlugin) AuthZReq(req authz.Request) authz.Response {

	if req.RequestBody != nil {
		var containerConfig dockerclient.ContainerConfig
		if err := json.NewDecoder(bytes.NewReader(req.RequestBody)).Decode(&containerConfig); err != nil {
			return authz.Response{Err: err.Error()}
		}

		//Container config may be empty - also sometime we may use another sturct?
		isAllowed := aclsAPI.ValidateRequest(req, containerConfig)

		if !isAllowed {
			return authz.Response{Msg: "Put some message here"}
		}
		return authz.Response{Allow: true}
	}
	return authz.Response{Msg: "Put some message here"}
}

//Before returning response to client
func (p *authzPlugin) AuthZRes(req authz.Request) authz.Response {
	return authz.Response{Allow: true}
}

/*
//An example with volumes...
func (p *authzPlugin) AuthZReq(req authz.Request) authz.Response {
	if req.RequestMethod == "POST" && startRegExp.MatchString(req.RequestURI) {
		// this is deprecated in docker, remove once hostConfig is dropped to
		// being available at start time
		if req.RequestBody != nil {
			type vfrom struct {
				VolumesFrom []string
			}
			vf := &vfrom{}
			if err := json.NewDecoder(bytes.NewReader(req.RequestBody)).Decode(vf); err != nil {
				return authz.Response{Err: err.Error()}
			}
			if len(vf.VolumesFrom) > 0 {
				goto noallow
			}
		}
		res := startRegExp.FindStringSubmatch(req.RequestURI)
		if len(res) < 1 {
			return authz.Response{Err: "unable to find container name"}
		}

		container, err := p.client.ContainerInspect(res[1])
		if err != nil {
			return authz.Response{Err: err.Error()}
		}
		image, _, err := p.client.ImageInspectWithRaw(container.Image, false)
		if err != nil {
			return authz.Response{Err: err.Error()}
		}
		if len(image.Config.Volumes) > 0 {
			goto noallow
		}
		for _, m := range container.Mounts {
			if m.Driver != "" {
				goto noallow
			}
		}
		if len(container.HostConfig.VolumesFrom) > 0 {
			goto noallow
		}
	}
	return authz.Response{Allow: true}

noallow:
	return authz.Response{Msg: "volumes are not allowed"}
}
*/
