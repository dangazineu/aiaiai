package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dangazineu/aiaiai/agent/pkg/domain"
	"github.com/labstack/echo/v4"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"log/slog"
	"net/http"
)

type TranslatorServer struct {
	logger        slog.Logger
	chatGptApiKey string
}

func NewTranslatorServer(logger slog.Logger, chatGptApiKey string) *TranslatorServer {
	return &TranslatorServer{
		logger:        logger,
		chatGptApiKey: chatGptApiKey,
	}
}

func (t TranslatorServer) Translate(c echo.Context) error {
	lexicalItem := c.Param("lexicalItem")
	if lexicalItem == "" {
		return c.String(http.StatusBadRequest, "lexical item wasn't provided")
	}

	translationResponse, err := t.translateByChatGPT(lexicalItem)
	if err != nil {
		t.logger.Error("Failed to translate", slog.Any("err", err))
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, translationResponse)
}

func (t TranslatorServer) translateByChatGPT(lexicalItem string) (*domain.TranslationResponse, error) {
	prompt := "Translate the lexical item, provide response in the following json format: lexicalItem(string), meaning (string), example (string). lexical item to translate:" + lexicalItem

	client := openai.NewClient(
		option.WithAPIKey(t.chatGptApiKey), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Model: openai.F(openai.ChatModelGPT3_5Turbo),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to make req to chat gpt api: %w", err)
	}
	println(chatCompletion.Choices[0].Message.Content)

	if len(chatCompletion.Choices) != 1 {
		return nil, fmt.Errorf("expected only one choice, but recieved: %d", len(chatCompletion.Choices))
	}

	var translationResp domain.TranslationResponse
	err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &translationResp)
	if err != nil {
		return nil, err
	}
	return &translationResp, nil
}
