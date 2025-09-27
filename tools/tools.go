// Package tools should expose all the tools available to agents
package tools

import (
	"encoding/json"

	"github.com/openai/openai-go/v2"
)

type Tooler interface {
	GetName() string
	Execute(input json.RawMessage) (string, error)
	GetFunctionStructure() openai.ChatCompletionToolUnionParam
}
