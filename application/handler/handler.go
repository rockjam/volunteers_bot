package handler

import (
	"dv/services"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type EventHandler struct {
	messageService services.Message
}

func NewEventHandler(messageService services.Message) EventHandler {
	return EventHandler{messageService}
}

func (e EventHandler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	err := e.messageService.HandleIncomingMessage([]byte(request.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}
