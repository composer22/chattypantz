package server

import (
	"fmt"
	"sync"
	"time"
)

var (
	maxChatRoomReq   = 1000                   // The maximum number of requests in the req channel.
	maxChatRoomSleep = 100 * time.Millisecond // How long to sleep between chan peeks.
)

// ChatRoom represents a hub of chatters where messages can be exchanged.
type ChatRoom struct {
	name     string            // The name of the room.
	chatters map[*Chatter]bool // A list of chatters in the room and if they are hidden from view.
	start    time.Time         // The start time of the room.
	lastReq  time.Time         // The last request time to the room.
	lastRsp  time.Time         // The last response time from the room.
	reqCount int64             // Total requests received.
	rspCount int64             // Total responses sent.
	reqq     chan *ChatRequest // Channel to receive requests.
	log      *ChatLogger       // Application log for events.
	wg       *sync.WaitGroup   // Wait group for the run from the chat room manager.
}

// ChatRoomNew is a factory function that returns a new instance of a chat room.
func ChatRoomNew(n string, cl *ChatLogger, g *sync.WaitGroup) *ChatRoom {
	return &ChatRoom{
		name:     n,
		chatters: make(map[*Chatter]bool),
		reqq:     make(chan *ChatRequest, maxChatRoomReq),
		log:      cl,
		wg:       g,
	}
}

// Run is the main routine that is evoked in background to accept commands to the room
func (r *ChatRoom) Run() {
	defer r.wg.Done()
	r.start = time.Now()
	for {
		select {
		case req, ok := <-r.reqq:
			if !ok { // Assume ch closed and shutdown notification
				return
			}
			r.lastReq = time.Now()
			r.reqCount++
			switch req.ReqType {
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
				r.sendResponse(req.Who, ChatRspTypeErrUnknownReq,
					fmt.Sprintf(`Unknown request sent to room "%s".`, r.name))
			}
		default:
			time.Sleep(maxChatRoomSleep)
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
	r.sendResponse(q.Who, ChatRspTypeListNames, fmt.Sprint(members))
}

// join adds the chatter to the room and notifies the group of the new chatters arrival.
func (r *ChatRoom) join(q *ChatRequest) {
	_, ok := r.chatters[q.Who]
	if ok {
		r.sendResponse(q.Who, ChatRspTypeErrAlreadyJoined,
			fmt.Sprintf(`You are already a member of room "%s".`, r.name))
		return
	}
	if q.Who.nickname == "" {
		r.sendResponse(q.Who, ChatRspTypeErrNicknameMandatory,
			fmt.Sprintf(`A nickname is mandatory to be a member of room "%s".`, r.name))
		return
	}
	r.sendResponseAll(ChatRspTypeJoin, fmt.Sprintf("%s has joined the room.", q.Who.nickname))
}

// hide hides/unhides a nickname from the user list
func (r *ChatRoom) hide(q *ChatRequest) {
	_, ok := r.chatters[q.Who]
	if !ok {
		r.sendResponse(q.Who, ChatRspTypeErrNotInRoom, fmt.Sprintf(`You are not a member of room "%s".`, r.name))
		return
	}
	r.chatters[q.Who] = !r.chatters[q.Who]
	htxt := "unhidden"
	t := ChatRspTypeUnhidden
	if r.chatters[q.Who] {
		htxt = "hidden"
		t = ChatRspTypeHidden
	}
	r.sendResponse(q.Who, t, fmt.Sprintf(`You are now %s in room "%s"`, htxt, r.name))
}

// message sends a message from a chatter to everyone in the room.
func (r *ChatRoom) message(q *ChatRequest) {
	_, ok := r.chatters[q.Who]
	if !ok {
		r.sendResponse(q.Who, ChatRspTypeErrNotInRoom, fmt.Sprintf(`You are not a member of room "%s".`, r.name))
		return
	}
	r.sendResponseAll(ChatRspTypeMessage, fmt.Sprintf("%s: %s", q.Who.nickname, q.Content))
}

// leave removes the chatter from the room and notifies the group the chatter has left.
func (r *ChatRoom) leave(q *ChatRequest) {
	_, ok := r.chatters[q.Who]
	if !ok {
		r.sendResponse(q.Who, ChatRspTypeErrNotInRoom, fmt.Sprintf(`You are not a member of room "%s".`, r.name))
		return
	}
	name := q.Who.nickname
	delete(r.chatters, q.Who)
	r.sendResponse(q.Who, ChatRspTypeLeave, fmt.Sprintf(`You have left room "%s".`, r.name))
	r.sendResponseAll(ChatRspTypeLeave, fmt.Sprintf("%s has left the room.", name))
}

// sendResponse sends a message to a single chatter in the room.
func (r *ChatRoom) sendResponse(c *Chatter, rt int, content string) {
	c.mu.Lock()
	if c.rspq != nil {
		r.lastRsp = time.Now()
		r.rspCount++
		c.rspq <- ChatResponseNew(r.name, rt, content)
	}
	c.mu.Unlock()
}

// sendResponseAll sends a message to all chatters in the room.
func (r *ChatRoom) sendResponseAll(rt int, content string) {

	rsp := ChatResponseNew(r.name, rt, content)
	for c := range r.chatters {
		c.mu.Lock()
		if c.rspq != nil {
			r.lastRsp = time.Now()
			r.rspCount++
			c.rspq <- rsp
		}
		c.mu.Unlock()
	}
}
