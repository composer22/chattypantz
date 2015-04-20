package server

import (
	"fmt"
	"sync"
)

// ChatRoom represents a hub of chatters where messages can be exchanged.
type ChatRoom struct {
	name     string            // The name of the room.
	chatters map[*Chatter]bool // A list of chatters in the room.
	reqq     chan *ChatRequest // Channel to receive commands.
	log      *ChatLogger       // Application log for events.
	wg       *sync.WaitGroup   // Wait group for the run.
}

// ChatRoomNew is a factory function that returns a new instance of a chat room.
func ChatRoomNew(n string, cl *ChatLogger, g *sync.WaitGroup) *ChatRoom {
	return &ChatRoom{
		name:     n,
		chatters: make(map[*Chatter]bool),
		reqq:     make(chan *ChatRequest),
		log:      cl,
		wg:       g,
	}
}

// Run is the main routine that is evoked in background to accept commands to the room
func (r *ChatRoom) Run() {
	defer r.wg.Done()
	for {
		select {
		case req, ok := <-r.reqq:
			if !ok { // Assume ch closed and shutdown notification
				return
			}
			switch req.reqType {
			case ChatReqTypeListNames:
				r.listNames(req)
			case ChatReqTypeJoin:
				r.join(req)
			case ChatReqTypeHide:
				r.hide(req)
			case ChatReqTypeMsg:
				r.message(req)
			case ChatReqTypeLeave:
				r.leave(req)
			default:
				r.sendResponse(req.who, ChatRspTypeErrUnknownReq, fmt.Sprintf(`Unknown request sent to room "%s".`, r.name))
			}
		}
	}
}

// listNames sends a response to the user with a list of all nicknames in the room.
func (r *ChatRoom) listNames(q *ChatRequest) {
	var members []string
	for c, hidden := range r.chatters {
		if !hidden { // don't return hidden names.
			members = append(members, c.nickname)
		}
	}
	r.sendResponse(q.who, ChatRspTypeListNames, fmt.Sprint(members))
}

// join adds the chatter to the room and notifies the group of the new chatters arrival.
func (r *ChatRoom) join(q *ChatRequest) {
	_, ok := r.chatters[q.who]
	if ok {
		r.sendResponse(q.who, ChatRspTypeErrAlreadyJoined, fmt.Sprintf(`You are already a member of room "%s".`, r.name))
		return
	}
	if q.who.nickname == "" {
		r.sendResponse(q.who, ChatRspTypeErrNicknameMandatory,
			fmt.Sprintf(`A nickname is mandatory to be a member of room "%s".`, r.name))
		return
	}
	r.sendResponseAll(ChatRspTypeJoin, fmt.Sprintf("%s has joined the room.", q.who.nickname))
}

// hide hides/unhides a nickname from the user list
func (r *ChatRoom) hide(q *ChatRequest) {
	_, ok := r.chatters[q.who]
	if !ok {
		r.sendResponse(q.who, ChatRspTypeErrNotInRoom, fmt.Sprintf(`You are not a member of room "%s".`, r.name))
		return
	}
	r.chatters[q.who] = !r.chatters[q.who]
	htxt := "unhidden"
	t := ChatRspTypeUnhidden
	if r.chatters[q.who] {
		htxt = "hidden"
		t = ChatRspTypeHidden
	}
	r.sendResponse(q.who, t, fmt.Sprintf(`You are now %s in room "%s"`, htxt, r.name))
}

// message sends a message from a chatter to everyone in the room.
func (r *ChatRoom) message(q *ChatRequest) {
	_, ok := r.chatters[q.who]
	if !ok {
		r.sendResponse(q.who, ChatRspTypeErrNotInRoom, fmt.Sprintf(`You are not a member of room "%s".`, r.name))
		return
	}
	r.sendResponseAll(ChatRspTypeMessage, fmt.Sprintf("%s: %s", q.who.nickname, q.content))
}

// leave removes the chatter from the room and notifies the group the chatter has left.
func (r *ChatRoom) leave(q *ChatRequest) {
	_, ok := r.chatters[q.who]
	if !ok {
		r.sendResponse(q.who, ChatRspTypeErrNotInRoom, fmt.Sprintf(`You are not a member of room "%s".`, r.name))
		return
	}
	name := q.who.nickname
	delete(r.chatters, q.who)
	r.sendResponse(q.who, ChatRspTypeLeave, fmt.Sprintf(`You have left room "%s".`, r.name))
	r.sendResponseAll(ChatRspTypeLeave, fmt.Sprintf("%s has left the room.", name))
}

// sendResponse sends a message to a single chatter in the room.
func (r *ChatRoom) sendResponse(c *Chatter, rt byte, content string) {
	c.mu.Lock()
	if c.rspq != nil {
		c.rspq <- ChatResponseNew(r.name, rt, content)
	}
	c.mu.Unlock()
}

// sendResponseAll sends a message to all chatters in the room.
func (r *ChatRoom) sendResponseAll(rt byte, content string) {
	rsp := ChatResponseNew(r.name, rt, content)
	for c := range r.chatters {
		c.mu.Lock()
		if c.rspq != nil {
			c.rspq <- rsp
		}
		c.mu.Unlock()
	}
}
