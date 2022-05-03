package proto

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type Template interface {
	Template() string
}

type TemplateRPC struct{ client *rpc.Client }

func (t *TemplateRPC) Template() string {
	var resp string

	_ = t.client.Call("Plugin.Template", new(interface{}), &resp)

	return resp
}

type TemplateRPCServer struct {
	Impl Template
}

func (t *TemplateRPCServer) Template(args interface{}, resp *string) error {
	*resp = t.Impl.Template()
	return nil
}

type TemplatePlugin struct {
	Impl Template
}

func (t *TemplatePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &TemplateRPCServer{Impl: t.Impl}, nil
}

func (TemplatePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &TemplateRPC{client: c}, nil
}
