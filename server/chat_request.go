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
	Who      *Chatter `json:"-"`        // The chatter who is issuing the request.
	RoomName string   `json:"roomName"` // The name of the room to receive the request.
	ReqType  int      `json:"reqType"`  // The command type ex: join, leave, send.
	Content  string   `json:"content"`  // Any message or text to interpret with the request.
}

// ChatMessageNew is a factory method that returns a new chat room message instance.
func ChatRequestNew(c *Chatter, m string, ct int, n string) *ChatRequest {
	return &ChatRequest{
		Who:      c,
		RoomName: m,
		ReqType:  ct,
		Content:  n,
	}
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (r *ChatRequest) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}
