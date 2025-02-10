package main

import (
	"github.com/dangazineu/aiaiai/whatsapp-agent/internal"
	"github.com/labstack/echo/v4"
	"log/slog"
	"os"
)

var (
	chatGptApiKey             = os.Getenv("OPENAI_API_KEY")
	whatsAppVerificationToken = os.Getenv("WHATSAPP_VERIFICATION_TOKEN")
	whatsAppApiToken          = os.Getenv("WHATSAPP_API_TOKEN")
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	e := echo.New()

	s := internal.NewServer(*logger, whatsAppVerificationToken, whatsAppApiToken, chatGptApiKey)
	e.GET("/", s.Subscribe)
	e.GET("/health", s.GetHealth)
	e.GET("/readiness", s.GetReadiness)
	e.POST("/", s.HandleMessage)

	logger.Error("server has failed", slog.Any("err", e.Start(":8080")))
}
