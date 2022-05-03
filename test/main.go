package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/pipego/plugin-score/test/proto"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "hclog",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command("./bin/plugin-score-template"),
		Logger:          logger,
	})
	defer client.Kill()

	rpcClient, _ := client.Client()

	raw, _ := rpcClient.Dispense("template")
	template := raw.(proto.Template)

	fmt.Println(template.Template())
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

var pluginMap = map[string]plugin.Plugin{
	"template": &proto.TemplatePlugin{},
}
