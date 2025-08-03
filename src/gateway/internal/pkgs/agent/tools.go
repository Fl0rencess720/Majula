package agent

import (
	"context"

	mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/viper"
)

func bingSearchTool(ctx context.Context) ([]tool.BaseTool, error) {
	cli, err := client.NewSSEMCPClient(viper.GetString("MCP_BING_SEARCH_URL"))
	if err != nil {
		return nil, err
	}
	if err := cli.Start(ctx); err != nil {
		return nil, err
	}
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "majula-client",
		Version: "1.0.0",
	}

	_, err = cli.Initialize(ctx, initRequest)
	if err != nil {
		return nil, err
	}

	tools, err := mcpp.GetTools(ctx, &mcpp.Config{Cli: cli})
	if err != nil {
		return nil, err
	}
	return tools, nil
}

func getCheckTools(ctx context.Context) ([]tool.BaseTool, error) {
	tools := []tool.BaseTool{}
	bingTools, err := bingSearchTool(ctx)
	if err != nil {
		return nil, err
	}
	return append(tools, bingTools...), nil
}

func toolsToInfo(ctx context.Context, tools []tool.BaseTool) ([]*schema.ToolInfo, error) {
	var toolsInfo []*schema.ToolInfo
	for _, t := range tools {
		info, err := t.Info(ctx)
		if err != nil {
			return nil, err
		}
		toolsInfo = append(toolsInfo, info)
	}
	return toolsInfo, nil
}
