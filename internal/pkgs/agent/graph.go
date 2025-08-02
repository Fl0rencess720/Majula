package agent

import (
	"context"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

type state struct {
	history []*schema.Message
}

func buildFactGraph(tpl prompt.ChatTemplate, cm model.ToolCallingChatModel) (*compose.Graph[map[string]any, *schema.Message], error) {
	compose.RegisterSerializableType[state]("state")
	g := compose.NewGraph[map[string]any, *schema.Message](
		compose.WithGenLocalState(func(ctx context.Context) *state {
			return &state{}
		}))
	err := g.AddChatTemplateNode(
		"ChatTemplate",
		tpl,
	)
	if err != nil {
		return nil, err
	}
	err = g.AddChatModelNode(
		"ChatModel",
		cm,
		compose.WithStatePreHandler(func(ctx context.Context, in []*schema.Message, state *state) ([]*schema.Message, error) {
			state.history = append(state.history, in...)
			return state.history, nil
		}),
		compose.WithStatePostHandler(func(ctx context.Context, out *schema.Message, state *state) (*schema.Message, error) {
			state.history = append(state.history, out)
			return out, nil
		}),
	)
	if err != nil {
		return nil, err
	}

	err = g.AddEdge(compose.START, "ChatTemplate")
	if err != nil {
		return nil, err
	}
	err = g.AddEdge("ChatTemplate", "ChatModel")
	if err != nil {
		return nil, err
	}
	err = g.AddEdge("ChatModel", compose.END)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func buildCheckGraph(ctx context.Context, tpl prompt.ChatTemplate, cm model.ToolCallingChatModel) (*compose.Graph[map[string]any, *schema.Message], error) {
	tools, err := getCheckTools(ctx)
	if err != nil {
		return nil, err
	}
	ragent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: cm,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: tools,
		},
	})
	if err != nil {
		return nil, err
	}
	compose.RegisterSerializableType[state]("state")
	g := compose.NewGraph[map[string]any, *schema.Message](
		compose.WithGenLocalState(func(ctx context.Context) *state {
			return &state{}
		}))

	err = g.AddChatTemplateNode(
		"ChatTemplate",
		tpl,
	)
	if err != nil {
		return nil, err
	}

	checkGraph, opt := ragent.ExportGraph()
	g.AddGraphNode("CheckNode", checkGraph, opt...)

	g.AddEdge(compose.START, "ChatTemplate")
	g.AddEdge("ChatTemplate", "CheckNode")
	g.AddEdge("CheckNode", compose.END)

	return g, nil

}
