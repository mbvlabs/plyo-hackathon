// Package tools should expose all the tools available to agents
package tools

import "encoding/json"

type Tooler interface {
	GetName() string
	GetDescription() string
	Execute(input json.RawMessage) (string, error)
}
