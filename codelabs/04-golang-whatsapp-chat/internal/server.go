package internal

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dangazineu/aiaiai/whatsapp-agent/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
)

type Server struct {
	logger              slog.Logger
	whatsAppApiToken    string
	whatsAppVeriyfToken string
	chatGptApiKey       string
}

func maskLeft(s string) string {
	rs := []rune(s)
	for i := 0; i < len(rs)-4; i++ {
		rs[i] = '*'
	}
	return string(rs)
}

func NewServer(logger slog.Logger, whatsAppVeriyfToken string, whatsAppApiToken string, chatGptApiKey string) *Server {
	logger.Info("Creating new server", slog.Any("whatsAppVerifyToken", maskLeft(whatsAppVeriyfToken)), slog.Any("whatsAppApiToken", maskLeft(whatsAppApiToken)), slog.Any("chatGptApiKey", maskLeft(chatGptApiKey)))
	return &Server{
		logger:              logger,
		whatsAppVeriyfToken: whatsAppVeriyfToken,
		whatsAppApiToken:    whatsAppApiToken,
		chatGptApiKey:       chatGptApiKey,
	}
}

func (s *Server) Subscribe(c echo.Context) error {
	mode := c.QueryParam("hub.mode")
	token := c.QueryParam("hub.verify_token")
	challenge := c.QueryParam("hub.challenge")

	if mode == "subscribe" && token == s.whatsAppVeriyfToken {
		return c.String(http.StatusOK, challenge)
	}
	return c.String(http.StatusForbidden, "invalid token")
}

func (s *Server) GetHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "healthy",
	})
}

func (s *Server) GetReadiness(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ready",
	})
}

func (s *Server) HandleMessage(c echo.Context) error {
	p := new(model.Payload)
	if err := c.Bind(p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	for _, entry := range p.Entry {
		for _, change := range entry.Changes {
			for _, message := range change.Value.Messages {
				if message.Text != nil {
					s.logger.Info("Received WhatsApp message", slog.Any("message", message))
					go s.respondToMessage(message)
				}
			}
		}
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

func (s *Server) respondToMessage(message model.Message) {
	client := openai.NewClient(
		option.WithAPIKey(s.chatGptApiKey), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a thoughtful comedian that always responds to everything with a joke that is related to the topic of the message you got. If you can't determine what is the topic of the conversation, make one up using the current time of the day as a starting point."),
			openai.UserMessage(message.Text.Body),
		}),
		Model: openai.F(openai.ChatModelGPT3_5Turbo),
	})
	if err != nil {
		s.logger.Error("failed to make req to chat gpt api", slog.Any("err", err))
		return
	}
	err = s.sendWhatsAppMessage(message.From, chatCompletion.Choices[0].Message.Content)
	if err != nil {
		s.logger.Error("failed to send message back to user", slog.Any("err", err))
		return
	}
	s.logger.Info("Responded to Message Successfully")
}

func (s *Server) sendWhatsAppMessage(to string, message string) error {
	requestBody := fmt.Sprintf(`{"messaging_product": "whatsapp", "recipient_type": "individual", "to": "%s", "type": "text", "text": {"preview_url": false, "body": "%s"}}`, to, message)

	req, err := http.NewRequest("POST", "https://graph.facebook.com/v21.0/490977760775365/messages", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error dumping request: %v", err))
	}
	s.logger.Info("Sending WhatsApp message", slog.Any("request", string(dump)))

	// Set the Authorization header after the log statement to avoid logging the token
	req.Header.Set("Authorization", "Bearer "+s.whatsAppApiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed post message to WhatsApp API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body from WhatsApp API: %s, %w", resp.Status, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message to WhatsApp API: %s, %s", resp.Status, string(body))
	}

	s.logger.Info("Sent WhatsApp message", slog.Any("response", string(body)))

	return nil
}
