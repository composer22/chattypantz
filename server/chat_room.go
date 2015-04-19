package server

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

var (
	maxChatRoomCmds  = 100                    // the maximum capacity of the cnd queue to the chat room.
	chatRoomMaxSleep = 250 * time.Millisecond // How long to sleep between msgq peeks.

)

// ChatRoom represents a hub of chatters where messages can be exchanged.
type ChatRoom struct {
	name     string            // The name of the room
	chatters map[*Chatter]bool // a list of chatters in the room
	Cmdq     chan string       // Channel to receive commands
	mu       sync.Mutex        // For locking access to chatter attributes.
	log      *ChatLogger       // Application log for events.

}

// ChatRoomNew is a factory function that returns a new instance of a chat room.
func ChatRoomNew(n string, cl *ChatLogger, addedOpts ...func(*ChatRoom)) *ChatRoom {
	cr := &ChatRoom{
		name:     n,
		chatters: make(map[*Chatter]bool),
		Cmdq:     make(chan string, maxChatRoomCmds),
		log:      cl,
	}
	// Additional hook for specialized custom options.
	for _, f := range addedOpts {
		f(cr)
	}
	return cr
}

// Run is the main routine that is evoked in background to accept commands to the room
func (r *ChatRoom) Run() {
Loop:
	for {
		select {
		case cmd, ok := <-r.Cmdq:
			if !ok { // Assume ch closed and shutdown notification
				break loop
			}
			// Do something with this command.

		default:
			time.Sleep(chatRoomMaxSleep) // Sleep before peeking again.
		}
	}
}

// join adds the chatter to the room and notifies the group of the new chatters arrival.
func (r *ChatRoom) join(c *Chatter) {
	defer r.wg.Lock()
	_, ok := r.chatters[c]
	if ok {
		r.sendMsg(c, "system", Sprintf(`You are already a member of room "%s".`, r.name))
		return
	}
	r.chatters[c] = true
	r.sendMsgAll(fmt.Sprintf("%s has joined the room.", c.Nickname))
}

// leave removes the chatter from the room and notifies the group the chatter has left.
func (r *ChatRoom) leave(c *Chatter) {
	defer r.wg.Lock()
	_, ok := r.chatters[c]
	if !ok {
		r.sendMsg(c, "system", Sprintf(`You are not a member of room "%s".`, r.name))
		return
	}
	name := c.Nickname
	delete(r.chatters[c])
	r.sendMsg(c, "system", Sprintf(`You have left room "%s".`, r.name))
	r.sendMsgAll(fmt.Sprintf("%s has left the room.", name))
}

// chatRoomMessage is a structure for JSON messages sent back to the client.
type ChatRoomMessage struct {
	RoomName string `json:"roomName"` // The room name where the message originated.
	MsgType  string `json:"msgType"`  // The message type ex: join, leave, send.
	Msg      string `json:"msg"`      // the message text
}

// String is an implentation of the Stringer interface so the structure is returned as a
// string to fmt.Print() etc.
func (m *chatRoomMessage) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

// sendMsg sends a message to a single member in the room.
func (r *ChatRoom) sendMsg(c *Chatter, mt string, msg string) {
	defer r.wg.Lock()
	crm := &ChatRoomMessage{
		RoomName: r.name,
		MsgType:  mt,
		Msg:      msg,
	}
	msgJSON, _ := fmt.Print(crm)
	c.Msgq <- msgJSON
}

// sendMsgAll sends a message to all members in the room.
func (r *ChatRoom) sendMsgAll(msg string) {
	defer r.wg.Lock()
	crm := &ChatRoomMessage{
		RoomName: r.name,
		MsgType:  "message",
		Msg:      msg,
	}
	msgJSON, _ := fmt.Print(crm)
	for c, _ := range r.chatters {
		c.Msgq <- msgJSON
	}
}
