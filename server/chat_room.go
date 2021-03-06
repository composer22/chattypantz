package server

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var (
	maxChatRoomReq = 1000 // The maximum number of requests in the req channel.
)

// ChatRoom represents a hub of chatters where messages can be exchanged.
type ChatRoom struct {
	mu       sync.RWMutex      // Lock against stats.
	name     string            // The name of the room.
	chatters map[*Chatter]bool // A list of chatters in the room and if they are hidden from view.
	start    time.Time         // The start time of the room.
	lastReq  time.Time         // The last request time to the room.
	lastRsp  time.Time         // The last response time from the room.
	reqCount uint64            // Total requests received.
	rspCount uint64            // Total responses sent.

	reqq chan *ChatRequest // Channel to receive requests.
	done chan bool         // Channel to receive signal to shutdown now.
	log  *ChatLogger       // Application log for events.
	wg   *sync.WaitGroup   // Wait group for the run from the chat room manager.
}

// ChatRoomNew is a factory function that returns a new instance of a chat room.
func ChatRoomNew(name string, d chan bool, cl *ChatLogger, g *sync.WaitGroup) *ChatRoom {
	return &ChatRoom{
		name:     name,
		chatters: make(map[*Chatter]bool),
		reqq:     make(chan *ChatRequest, maxChatRoomReq),
		done:     d,
		log:      cl,
		wg:       g,
	}
}

// Run is the main routine that is evoked in background to accept commands to the room.
func (r *ChatRoom) Run() {
	defer r.wg.Done()
	r.start = time.Now()
	for {
		select {
		case <-r.done: // Server signal quit
			return
		case req, ok := <-r.reqq:
			if !ok { // Assume ch closed and shutdown notification
				return
			}
			r.mu.Lock()
			r.lastReq = time.Now()
			r.reqCount++
			r.mu.Unlock()
			switch req.ReqType {
			case ChatReqTypeListNames:
				r.listNames(req)
			case ChatReqTypeJoin:
				r.join(req)
			case ChatReqTypeHide:
				r.hide(req)
			case ChatReqTypeUnhide:
				r.unhide(req)
			case ChatReqTypeMsg:
				r.message(req)
			case ChatReqTypeLeave:
				r.leave(req)
			default:
				r.sendResponse(req.Who, ChatRspTypeErrUnknownReq,
					fmt.Sprintf(`Unknown request sent to room "%s".`, r.Name()), nil)
			}
		}
		runtime.Gosched()
	}
}

// join adds the chatter to the room and notifies the group of the new chatters arrival.
func (r *ChatRoom) join(q *ChatRequest) {
	switch {
	case r.isMember(q.Who):
		r.sendResponse(q.Who, ChatRspTypeErrAlreadyJoined,
			fmt.Sprintf(`You are already a member of room "%s".`, r.Name()), nil)
	case r.isMemberName(q.Who.Nickname()):
		r.sendResponse(q.Who, ChatRspTypeErrNicknameUsed,
			fmt.Sprintf(`Nickname "%s" is already in use in room "%s".`, q.Who.Nickname(), r.Name()), nil)
	default:
		r.mu.Lock()
		r.chatters[q.Who] = false
		if q.Content == "hidden" {
			r.chatters[q.Who] = true
		}
		var names []string
		for c, hidden := range r.chatters {
			if !hidden { // don't return hidden names.
				names = append(names, c.Nickname())
			}
		}
		r.mu.Unlock()
		r.sendResponseAll(ChatRspTypeJoin, fmt.Sprintf("%s has joined the room.", q.Who.Nickname()), names)
	}
}

// listNames sends a response to the user with a list of all nicknames in the room.
func (r *ChatRoom) listNames(q *ChatRequest) {
	var names []string
	r.mu.RLock()
	for c, hidden := range r.chatters {
		if !hidden { // don't return hidden names.
			names = append(names, c.Nickname())
		}
	}
	r.mu.RUnlock()
	r.sendResponse(q.Who, ChatRspTypeListNames, "", names)
}

// hide visually makes a nickname inactive in the user list
func (r *ChatRoom) hide(q *ChatRequest) {
	r.mu.Lock()
	r.chatters[q.Who] = true
	r.mu.Unlock()
	r.sendResponse(q.Who, ChatRspTypeHide, fmt.Sprintf(`You are now hidden in room "%s".`, r.Name()), nil)
}

// unhide visually makes a nickname active in the user list
func (r *ChatRoom) unhide(q *ChatRequest) {
	r.mu.Lock()
	r.chatters[q.Who] = false
	r.mu.Unlock()
	r.sendResponse(q.Who, ChatRspTypeUnhide, fmt.Sprintf(`You are now unhidden in room "%s".`, r.Name()), nil)
}

// message sends a message from a chatter to everyone in the room.
func (r *ChatRoom) message(q *ChatRequest) {
	r.mu.RLock()
	isHidden := r.chatters[q.Who]
	r.mu.RUnlock()
	if isHidden {
		r.sendResponse(q.Who, ChatRspTypeErrHiddenNickname,
			fmt.Sprintf(`Nickname "%s" is hidden. Cannot post in room "%s".`, q.Who.Nickname(),
				r.name), nil)
	} else {
		r.sendResponseAll(ChatRspTypeMsg, fmt.Sprintf("%s: %s", q.Who.Nickname(), q.Content), nil)
	}
}

// leave removes the chatter from the room and notifies the group the chatter has left.
func (r *ChatRoom) leave(q *ChatRequest) {
	if ok := r.isMember(q.Who); !ok {
		return
	}
	name := q.Who.Nickname()
	var names []string
	r.mu.Lock()
	delete(r.chatters, q.Who)
	for c, hidden := range r.chatters {
		if !hidden { // don't return hidden names.
			names = append(names, c.Nickname())
		}
	}
	r.mu.Unlock()
	r.sendResponse(q.Who, ChatRspTypeLeave, fmt.Sprintf(`You have left room "%s".`, r.Name()), nil)
	r.sendResponseAll(ChatRspTypeLeave, fmt.Sprintf("%s has left the room.", name), names)
}

// ChatRoomStats is a simple structure for returning statistic information on the room.
type ChatRoomStats struct {
	Name     string                 `json:"name"`     // The name of the room.
	Start    time.Time              `json:"start"`    // The start time of the room.
	LastReq  time.Time              `json:"lastReq"`  // The last request time to the room.
	LastRsp  time.Time              `json:"lastRsp"`  // The last response time from the room.
	ReqCount uint64                 `json:"reqcount"` // Total requests received.
	RspCount uint64                 `json:"rspCount"` // Total responses sent.
	Chatters []*ChatRoomChatterStat `json:"chatters"` // Stats on chatters in the room
}

type ChatRoomChatterStat struct {
	Nickname   string `json:"nickname"`   // The nickname of the chatter.
	RemoteAddr string `json:"remoteAddr"` // The remote IP and port of the chatter.
}

// ChatRoomStatsNew returns status information on the room.
func (r *ChatRoom) ChatRoomStatsNew() *ChatRoomStats {
	r.mu.RLock()
	defer r.mu.RUnlock()
	stat := &ChatRoomStats{
		Name:     r.name,
		Start:    r.start,
		LastReq:  r.lastReq,
		LastRsp:  r.lastRsp,
		ReqCount: r.reqCount,
		RspCount: r.rspCount,
		Chatters: []*ChatRoomChatterStat{},
	}
	for ctr := range r.chatters {
		ctrStat := ctr.ChatterStatsNew()
		stat.Chatters = append(stat.Chatters, &ChatRoomChatterStat{
			Nickname:   ctrStat.Nickname,
			RemoteAddr: ctrStat.RemoteAddr,
		})
	}
	return stat
}

// isEmpty validates whether the room is empty of chatters.
func (r *ChatRoom) isEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.chatters) > 0 {
		return false
	}
	return true
}

// isMember validates if the member exists in the room.
func (r *ChatRoom) isMember(c *Chatter) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.chatters[c]
	return ok
}

// isMemberName validates if a member is using a nickname in the room.
func (r *ChatRoom) isMemberName(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for c := range r.chatters {
		if c.Nickname() == name {
			return true
		}
	}
	return false
}

// sendResponse sends a message to a single chatter in the room.
func (r *ChatRoom) sendResponse(c *Chatter, rspt int, cont string, l []string) {
	if l == nil {
		l = []string{}
	}
	c.sendResponse(r.Name(), rspt, cont, l)
	r.mu.Lock()
	r.lastRsp = time.Now()
	r.rspCount++
	r.mu.Unlock()
}

// sendResponseAll sends a message to all chatters in the room.
func (r *ChatRoom) sendResponseAll(rspt int, cont string, l []string) {
	if l == nil {
		l = []string{}
	}
	r.mu.Lock()
	for c := range r.chatters {
		c.sendResponse(r.name, rspt, cont, l)
		r.lastRsp = time.Now()
		r.rspCount++
	}
	r.mu.Unlock()
}

// Name returns the current name of the room.
func (r *ChatRoom) Name() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.name
}

// SetName sets the name of the room.
func (r *ChatRoom) SetName(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.name = name
}
