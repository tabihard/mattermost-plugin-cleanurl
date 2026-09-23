package main

import (
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost/server/public/plugin"
)

// HelloWorldPlugin replies to any request under /plugins/<plugin id> with a
// simple greeting, following the Mattermost "Hello, World!" server plugin
// tutorial: https://developers.mattermost.com/integrate/plugins/components/server/hello-world/
type HelloWorldPlugin struct {
	plugin.MattermostPlugin
}

func (p *HelloWorldPlugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func main() {
	plugin.ClientMain(&HelloWorldPlugin{})
}
