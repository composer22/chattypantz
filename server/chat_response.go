package server

import "encoding/json"

const (
	ChatRspTypeList = iota
	ChatRspTypeJoin
	ChatRspTypeAlreadyJoined
	ChatRspTypeHidden
	ChatRspTypeUnhidden
	ChatRspTypeChat
	ChatRspTypeLeave
	ChatRspTypeNotInRoom
	ChatRspTypeUnknownCmd
)

// ChatResponse is a structure for JSON responses sent back to the client.
type ChatResponse struct {
	RoomName string `json:"roomName"` // The room name where the response originated.
	RpsType  int    `json:"rpsType"`  // The response type ex: join, leave, send.
	Msg      string `json:"msg"`      // the message text.
}

// ChatResponseNew is a factory method that returns a new chat room message instance.
func ChatResponseNew(n string, rt int, msg string, addedOpts ...func(*ChatResponse)) *ChatResponse {
	rsp := &ChatResponse{
		RoomName: n,
		RpsType:  rt,
		Msg:      msg,
	}
	// Additional hook for specialized custom options.
	for _, f := range addedOpts {
		f(rsp)
	}
	return rsp
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (c *ChatResponse) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}
