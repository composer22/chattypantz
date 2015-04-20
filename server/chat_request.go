package server

import "encoding/json"

const (
	ChatReqTypeSetNickname = iota
	ChatReqTypeGetNickname
	ChatReqTypeListRooms
	ChatReqTypeListNames
	ChatReqTypeJoin
	ChatReqTypeHide
	ChatReqTypeMsg
	ChatReqTypeLeave
)

// ChatRequest is a structure for commands sent for processing from the client.
type ChatRequest struct {
	who      *Chatter `json:"-"`        // The chatter who is issuing the request.
	roomName string   `json:"roomName"` // The name of the room to receive the request.
	reqType  int      `json:"reqType"`  // The command type ex: join, leave, send.
	content  string   `json:"content"`  // Any message or text to interpret with the request.
}

// ChatMessageNew is a factory method that returns a new chat room message instance.
func ChatRequestNew(c *Chatter, m string, ct int, n string) *ChatRequest {
	return &ChatRequest{
		who:      c,
		roomName: m,
		reqType:  ct,
		content:  n,
	}
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (r *ChatRequest) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}
