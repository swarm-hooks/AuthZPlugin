// authz project AuthZPluginBackEnd.go
package authz

import (
	"github.com/AuthZPluginBackEnd/handlers"
	"net/http"

	"github.com/docker/docker/pkg/authorization"
)

const (
	manifest = `{"Implements": ["` + authorization.AuthZApiImplements + `"]}`
	reqPath  = "/" + authorization.AuthZApiRequest
	resPath  = "/" + authorization.AuthZApiResponse
)

// Request is the structure that docker's requests are deserialized to.
type Request authorization.Request

// Response is the strucutre that the plugin's responses are serialized to.
type Response authorization.Response

// Plugin represent the interface a plugin must fulfill.
type Plugin interface {
	AuthZReq(Request) Response
	AuthZRes(Request) Response
}

type Handler struct {
	plugin Plugin
	handlers.Handler
}

// NewHandler initializes the request handler with a plugin implementation.
func NewHandler(plugin Plugin) *Handler {
	h := &Handler{plugin, handlers.NewHandler(manifest)}
	h.initMux()
	return h
}

func (h *Handler) initMux() {
	h.handle(reqPath, func(req Request) Response {
		return h.plugin.AuthZReq(req)
	})

	h.handle(resPath, func(req Request) Response {
		return h.plugin.AuthZRes(req)
	})
}

type actionHandler func(Request) Response

func (h *Handler) handle(name string, actionCall actionHandler) {
	h.HandleFunc(name, func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := handlers.DecodeRequest(w, r, &req); err != nil {
			return
		}

		res := actionCall(req)

		handlers.EncodeResponse(w, res, res.Err)
	})
}
