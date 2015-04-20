package server

import (
	"fmt"
	"sync"
	"time"
)

var (
	maxChatRoomCmds  = 1000                   // the maximum capacity of the cnd queue to the chat room.
	maxChatRoomSleep = 100 * time.Millisecond // How long to sleep between msgq peeks.
)

// ChatRoom represents a hub of chatters where messages can be exchanged.
type ChatRoom struct {
	name     string            // The name of the room
	chatters map[*Chatter]bool // a list of chatters in the room
	Cmdq     chan *ChatCommand // Channel to receive commands
	mu       sync.Mutex        // For locking access to chatter attributes.
	log      *ChatLogger       // Application log for events.
}

// ChatRoomNew is a factory function that returns a new instance of a chat room.
func ChatRoomNew(n string, cl *ChatLogger, addedOpts ...func(*ChatRoom)) *ChatRoom {
	cr := &ChatRoom{
		name:     n,
		chatters: make(map[*Chatter]bool),
		Cmdq:     make(chan *ChatCommand, maxChatRoomCmds),
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
loop:
	for {
		select {
		case cmd, ok := <-r.Cmdq:
			if !ok { // Assume ch closed and shutdown notification
				break loop
			}
			// Respond to the received command
			switch cmd.cmdType {
			case ChatCmdTypeList:
				r.list(cmd)
			case ChatCmdTypeJoin:
				r.join(cmd)
			case ChatCmdTypeHide:
				r.hide(cmd)
			case ChatCmdTypeUnhide:
				r.unhide(cmd)
			case ChatCmdTypeChat:
				r.chat(cmd)
			case ChatCmdTypeLeave:
				r.leave(cmd)
			default:
				r.sendResponse(cmd.who, ChatRspTypeUnknownCmd, fmt.Sprintf(`Unknown command sent to room "%s".`, r.name))
			}
		default:
			time.Sleep(maxChatRoomSleep) // Sleep before peeking again.
		}
	}
}

// list sends a message to the user of everyone's name in the room.
func (r *ChatRoom) list(c *ChatCommand) {
	defer r.mu.Lock()
	var members []string
	for ct, hidden := range r.chatters {
		if !hidden {
			members = append(members, ct.nickname)
		}
	}
	r.sendResponse(c.who, ChatRspTypeList, fmt.Sprint(members))
}

// join adds the chatter to the room and notifies the group of the new chatters arrival.
func (r *ChatRoom) join(c *ChatCommand) {
	defer r.mu.Lock()
	_, ok := r.chatters[c.who]
	if ok {
		r.sendResponse(c.who, ChatRspTypeAlreadyJoined, fmt.Sprintf(`You are already a member of room "%s".`, r.name))
		return
	}
	r.chatters[c.who] = false
	r.sendResponseAll(ChatRspTypeJoin, fmt.Sprintf("%s has joined the room.", c.who.nickname))
}

// hide hides a user from the user list
func (r *ChatRoom) hide(c *ChatCommand) {
	defer r.mu.Lock()
	_, ok := r.chatters[c.who]
	if !ok {
		r.sendResponse(c.who, ChatRspTypeNotInRoom, fmt.Sprintf(`You are not a member of room "%s".`, r.name))
		return
	}
	r.chatters[c.who] = true
	r.sendResponse(c.who, ChatRspTypeHidden, fmt.Sprintf(`You are now hidden in room "%s"`, r.name))
}

// unhide makes a user visible in the user list
func (r *ChatRoom) unhide(c *ChatCommand) {
	defer r.mu.Lock()
	_, ok := r.chatters[c.who]
	if !ok {
		r.sendResponse(c.who, ChatRspTypeNotInRoom, fmt.Sprintf(`You are not a member of room "%s".`, r.name))
		return
	}
	r.chatters[c.who] = false
	r.sendResponse(c.who, ChatRspTypeUnhidden, fmt.Sprintf(`You are now unhidden in room "%s"`, r.name))
}

// chat sends a message from a user to everyone in the room.
func (r *ChatRoom) chat(c *ChatCommand) {
	defer r.mu.Lock()
	_, ok := r.chatters[c.who]
	if !ok {
		r.sendResponse(c.who, ChatRspTypeNotInRoom, fmt.Sprintf(`You are not a member of room "%s".`, r.name))
		return
	}
	r.sendResponseAll(ChatRspTypeChat, fmt.Sprintf("%s: %s", c.who.nickname, c.msg))
}

// leave removes the chatter from the room and notifies the group the chatter has left.
func (r *ChatRoom) leave(c *ChatCommand) {
	defer r.mu.Lock()
	_, ok := r.chatters[c.who]
	if !ok {
		r.sendResponse(c.who, ChatRspTypeNotInRoom, fmt.Sprintf(`You are not a member of room "%s".`, r.name))
		return
	}
	name := c.who.nickname
	delete(r.chatters, c.who)
	r.sendResponse(c.who, ChatRspTypeLeave, fmt.Sprintf(`You have left room "%s".`, r.name))
	r.sendResponseAll(ChatRspTypeLeave, fmt.Sprintf("%s has left the room.", name))
}

// sendResponse sends a message to a single member in the room.
func (r *ChatRoom) sendResponse(c *Chatter, rt int, msg string) {
	defer r.mu.Lock()
	if c.rspq != nil {
		c.rspq <- ChatResponseNew(r.name, rt, msg)
	}
}

// sendResponseAll sends a message to all members in the room.
func (r *ChatRoom) sendResponseAll(rt int, msg string) {
	defer r.mu.Lock()
	rsp := ChatResponseNew(r.name, rt, msg)
	for c := range r.chatters {
		if c.rspq != nil {
			c.rspq <- rsp
		}
	}
}
