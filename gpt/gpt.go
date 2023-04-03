package gpt

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type ChatBGIResponseBody struct {
	Type            string `json:"type"`
	Message         string `json:"message"`
	Conversation_ID string `json:"conversation_id"`
	Parent_ID       string `json:"parent_id"`
	Use_Paid        bool   `json:"use_paid"`
}

type ChatBGIRequestBody struct {
	Message         string `json:"message"`
	New_Title       string `json:"new_title"`
	Use_Paid        bool   `json:"use_paid"`
	Conversation_ID string `json:"conversation_id,omitempty"`
}

func Completions(msg string, conversation_id string) (string, string, error) {
	var reply string
	var new_conversation_id string
	var resErr error

	for retry := 1; retry <= 3; retry++ {
		if retry > 1 {
			time.Sleep(time.Duration(retry-1) * 100 * time.Millisecond)
		}
		reply, new_conversation_id, resErr = websocketCompletions(msg, conversation_id, retry)
		log.Printf("gpt request(%d) json: %s\n", retry, reply)

		if resErr == nil {
			break
		}
		log.Printf("gpt request(%d) error: %v\n", retry, resErr)
	}
	if resErr != nil {
		return "", "", resErr
	}
	return reply, new_conversation_id, nil
}

func websocketCompletions(msg string, conversation_id string, runtimes int) (string, string, error) {

	RequestBody := ChatBGIRequestBody{
		Message:         msg,
		New_Title:       "test",
		Use_Paid:        true,
		Conversation_ID: conversation_id,
	}
	requestData, err := json.Marshal(RequestBody)
	if err != nil {
		return "", "", fmt.Errorf("json.Marshal requestBody error: %v", err)
	}

	log.Printf("gpt request(%d) json: %s\n", runtimes, string(requestData))

	u := url.URL{Scheme: "ws", Host: "localhost:8000", Path: "/conv"}

	req := http.Request{
		Header: http.Header{},
	}
	req.Header.Set("Cookie", "user_auth=wechat")

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), req.Header)

	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, requestData)

	if err != nil {
		return "", "", fmt.Errorf("write message error: %v", err)
	}

	var result, new_conversation_id string
	for {
		_, jsonMessage, err := conn.ReadMessage()

		if err != nil {
			break
		}

		var receivedMessage ChatBGIResponseBody

		err = json.Unmarshal(jsonMessage, &receivedMessage)
		if err != nil {
			log.Println("json:", err)
			return "", "", fmt.Errorf("load json error: %v", err)
		}
		result = receivedMessage.Message
		new_conversation_id = receivedMessage.Conversation_ID
	}
	return result, new_conversation_id, nil
}
