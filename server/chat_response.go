package server

import (
	"encoding/json"
	"errors"
)

const (
	ChatRspTypeSetNickname = 101 + iota
	ChatRspTypeGetNickname
	ChatRspTypeListRooms
	ChatRspTypeJoin
	ChatRspTypeListNames
	ChatRspTypeHidden
	ChatRspTypeUnhidden
	ChatRspTypeMessage
	ChatRspTypeLeave
)

const (
	ChatRspTypeErrRoomMandatory = 1001 + iota
	ChatRspTypeErrMaxRoomsReached
	ChatRspTypeErrNicknameMandatory
	ChatRspTypeErrAlreadyJoined
	ChatRspTypeErrNicknameUsed
	ChatRspTypeErrNotInRoom
	ChatRspTypeErrUnknownReq
)

// ChatResponse is a structure for JSON responses sent back to the client.
type ChatResponse struct {
	RoomName string   `json:"roomName"` // The room name where the response originated.
	RspType  int      `json:"rspType"`  // The response type ex: join, leave, send.
	Content  string   `json:"content"`  // Any message text or other content for the client.
	List     []string `json:"list"`     // A list of entries returned with the response.
}

// ChatResponseNew is a factory method that returns a new chat room message instance.
func ChatResponseNew(m string, rt int, c string, l []string) (*ChatResponse, error) {
	if rt < ChatRspTypeSetNickname ||
		(rt > ChatRspTypeLeave && rt < ChatRspTypeErrRoomMandatory) ||
		rt > ChatRspTypeErrUnknownReq {
		return nil, errors.New("Response Type is out of range.")
	}
	return &ChatResponse{
		RoomName: m,
		RspType:  rt,
		Content:  c,
		List:     l,
	}, nil
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (r *ChatResponse) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}
