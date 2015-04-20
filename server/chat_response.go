package server

import "encoding/json"

const (
	ChatRspTypeSetNickname = iota
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
	RspType  int    `json:"rspType"`  // The response type ex: join, leave, send.
	Content  string `json:"content"`  // Any message text or other content for the client.
}

// ChatResponseNew is a factory method that returns a new chat room message instance.
func ChatResponseNew(m string, rt int, c string) *ChatResponse {
	return &ChatResponse{
		RoomName: m,
		RspType:  rt,
		Content:  c,
	}
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (r *ChatResponse) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}
