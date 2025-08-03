package agent

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	openai2 "github.com/cloudwego/eino-ext/libs/acl/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/spf13/viper"
)

func newChatModel(ctx context.Context) (model.ToolCallingChatModel, error) {
	cm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey: viper.GetString("CHATMODEL_API_KEY"),
		Model:  viper.GetString("chatmodel.model"),
		ResponseFormat: &openai2.ChatCompletionResponseFormat{
			Type: openai2.ChatCompletionResponseFormatTypeJSONObject,
		},
		BaseURL: viper.GetString("chatmodel.baseURL"),
	})

	if err != nil {
		return nil, err
	}
	return cm, nil
}
