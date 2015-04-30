package server

import (
	"encoding/json"
	"errors"
)

const (
	ChatReqTypeSetNickname = 101 + iota
	ChatReqTypeGetNickname
	ChatReqTypeListRooms
	ChatReqTypeJoin
	ChatReqTypeListNames
	ChatReqTypeHide
	ChatReqTypeUnhide
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
func ChatRequestNew(c *Chatter, room string, reqt int, cont string) (*ChatRequest, error) {
	if reqt < ChatReqTypeSetNickname || reqt > ChatReqTypeLeave {
		return nil, errors.New("Request Type is out of range.")
	}
	return &ChatRequest{
		Who:      c,
		RoomName: room,
		ReqType:  reqt,
		Content:  cont,
	}, nil
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (r *ChatRequest) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}
