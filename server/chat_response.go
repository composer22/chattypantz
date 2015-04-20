package server

import "encoding/json"

const (
	ChatRspTypeSetNickname byte = iota
	ChatRspTypeNickname
	ChatRspTypeListRooms
	ChatRspTypeListNames
	ChatRspTypeJoin
	ChatRspTypeErrRoomIsMandatory
	ChatRspTypeErrMaxRoomsReached
	ChatRspTypeErrNicknameMandatory
	ChatRspTypeErrAlreadyJoined
	ChatRspTypeHidden
	ChatRspTypeUnhidden
	ChatRspTypeMessage
	ChatRspTypeLeave
	ChatRspTypeErrNotInRoom
	ChatRspTypeErrUnknownReq
)

// ChatResponse is a structure for JSON responses sent back to the client.
type ChatResponse struct {
	RoomName string `json:"roomName"` // The room name where the response originated.
	RspType  byte   `json:"rspType"`  // The response type ex: join, leave, send.
	Content  string `json:"content"`  // The message text.
}

// ChatResponseNew is a factory method that returns a new chat room message instance.
func ChatResponseNew(m string, rt byte, c string) *ChatResponse {
	return &ChatResponse{
		RoomName: m,
		RspType:  rt,
		Content:  c,
	}
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (c *ChatResponse) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}
