package server

import "encoding/json"

const (
	ChatCmdTypeList = iota
	ChatCmdTypeJoin
	ChatCmdTypeHide
	ChatCmdTypeUnhide
	ChatCmdTypeChat
	ChatCmdTypeLeave
)

// ChatCommand is a structure for commands sent to the room.
type ChatCommand struct {
	who     *Chatter `json:"chatter"` // The chatter who is issuing the command.
	cmdType int      `json:"cmdType"` // The command type ex: join, leave, send.
	msg     string   `json:"msg"`     // Any message or additional text to interpret with the command
}

// ChatMessageNew is a factory method that returns a new chat room message instance.
func ChatCommandNew(c *Chatter, ct int, msg string, addedOpts ...func(*ChatCommand)) *ChatCommand {
	cmd := &ChatCommand{
		who:     c,
		cmdType: ct,
		msg:     msg,
	}
	// Additional hook for specialized custom options.
	for _, f := range addedOpts {
		f(cmd)
	}
	return cmd
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (m *ChatCommand) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}
