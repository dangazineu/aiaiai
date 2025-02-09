package main

import (
	"github.com/dangazineu/aiaiai/agent/pkg/server"
	"github.com/labstack/echo/v4"
	"log/slog"
	"os"
)

var (
	chatGptApiKey = os.Getenv("OPENAI_API_KEY")
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	e := echo.New()

	translatorServer := server.NewTranslatorServer(*logger, chatGptApiKey)
	e.GET("/v1/translations/:lexicalItem", translatorServer.Translate)

	logger.Error("server has failed", slog.Any("err", e.Start(":8080")))
}
