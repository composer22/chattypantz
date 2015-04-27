package server

import (
	"fmt"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

var (
	maxChatterRsp = 1000 // The max number of responses in the response channel.
)

// Chatter is a wrapper around a connection that represents one chat client on the server.
// It is a parasitic class in that it lives and dies within the context of the server and is trusted
// to use and modify server atttributes directly.
type Chatter struct {
	cMngr    *ChatManager       // The chat manager this chatter is attached to.
	ws       *websocket.Conn    // The socket to the remote client.
	nickname string             // The friendly nickname to display in a conversation.
	maxIdle  int                // Maximum idle time a ws should optionally wait.
	start    time.Time          // The start time of the connection.
	lastReq  time.Time          // The last request time of the connection.
	lastRsp  time.Time          // The last response time to the connection.
	reqCount uint64             // Total requests received.
	rspCount uint64             // Total responses sent.
	done     chan bool          // Signal that chatter is closed.
	rspq     chan *ChatResponse // A channel to receive information to send to the remote client.
	log      *ChatLogger        // Server logger
	mu       sync.Mutex         // For locking access to chatter attributes.
	wg       sync.WaitGroup     // Synchronization of channel close.

}

// ChatterNew is a factory function that returns a new Chatter instance
func ChatterNew(cm *ChatManager, c *websocket.Conn, mi int, l *ChatLogger) *Chatter {
	return &Chatter{
		cMngr:   cm,
		ws:      c,
		maxIdle: mi,
		done:    make(chan bool, 1),
		rspq:    make(chan *ChatResponse, maxChatterRsp),
	}
}

// Run starts the event loop that manages the sending and receiving of information to the client.
func (c *Chatter) Run() {
	c.start = time.Now()
	c.cMngr.wg.Add(1) // We let the big boss also perform waits for chatters, so it can close down,
	c.wg.Add(1)       //   but we also have our own in send().
	go c.send()       // Spawn response handling to the client in the background.
	c.receive()       // Then wait on incoming requests.
}

// receive polls and handles any commands or information sent from the remote client.
func (c *Chatter) receive() {
	defer c.cMngr.wg.Done()
	remoteAddr := c.ws.Request().RemoteAddr
	for {
		// Set optional idle timeout on the receive.
		if c.maxIdle > 0 {
			c.ws.SetReadDeadline(time.Now().Add(time.Duration(c.maxIdle) * time.Second))
		}
		var req ChatRequest
		if err := websocket.JSON.Receive(c.ws, &req); err != nil {
			e, ok := err.(net.Error)
			switch {
			case ok && e.Timeout():
				c.log.LogSession("disconnected", remoteAddr, "Client forced to disconnect due to inactivity.")
				c.shutDown()
			case err.Error() == "EOF":
				c.log.LogSession("disconnected", remoteAddr, "Client disconnected.")
				c.shutDown()
			case strings.Contains(err.Error(), "use of closed network connection"): // cntl-c safety.
				c.shutDown()
			default:
				c.log.LogError(remoteAddr, fmt.Sprintf("Couldn't receive. Error: %s", err.Error()))
			}
			return
		}
		c.mu.Lock()
		c.lastReq = time.Now()
		c.reqCount++
		c.mu.Unlock()
		c.log.LogSession("received", remoteAddr, fmt.Sprintf("%s", &req))
		switch req.ReqType {
		case ChatReqTypeSetNickname:
			c.setNickname(&req)
		case ChatReqTypeGetNickname:
			c.getNickname()
		case ChatReqTypeListRooms:
			c.listRooms()
		default: // let room handle other requests or send error if no room name provided.
			req.Who = c
			c.sendRequestToRoom(&req)
		}
	}
}

// send is a go routine used to poll queued messages to send to the client.
func (c *Chatter) send() {
	defer c.wg.Done()
	remoteAddr := fmt.Sprint(c.ws.Request().RemoteAddr)
	for {
		select {
		case <-c.cMngr.done: // Server shutdown signal.
			c.ws.Close() // Break the receive() loop and force a chatter shutdown.
			return
		case <-c.done: // Chatter shutdown signal.
			c.ws.Close() // Break the receive() loop and force a chatter shutdown.
			return
		case rsp, ok := <-c.rspq:
			if !ok { // Assume ch closed might also be shutdown notification from somebody.
				c.ws.Close() // Break the receive() looper and force a chatter shutdown.
				return
			}
			c.mu.Lock()
			c.lastRsp = time.Now()
			c.rspCount++
			c.mu.Unlock()
			c.log.LogSession("sent", remoteAddr, fmt.Sprintf("%s", rsp))
			if err := websocket.JSON.Send(c.ws, rsp); err != nil {
				switch {
				case err.Error() == "EOF":
					c.log.LogSession("disconnected", remoteAddr, "Client disconnected.")
					return
				default:
					c.log.LogError(remoteAddr, fmt.Sprintf("Couldn't send. Error: %s", err.Error()))
				}
			}
		}
		runtime.Gosched()
	}
}

// shutDown shuts down sending/receiving.
func (c *Chatter) shutDown() {
	close(c.done) // Signal to send() and rooms we are quitting.
	c.wg.Wait()   // Wait for send()
	c.cMngr.removeChatterAllRooms(c)
}

// setNickname sets the nickname for the chatter.
func (c *Chatter) setNickname(r *ChatRequest) {
	if r.Content == "" {
		c.sendResponse(ChatRspTypeErrNicknameMandatory, "Nickname cannot be blank.", nil)
		return
	}
	c.mu.Lock()
	c.nickname = r.Content
	c.mu.Unlock()
	c.sendResponse(ChatRspTypeSetNickname, fmt.Sprintf(`Nickname set to "%s".`, c.Nickname()), nil)
}

// nickname returns the nickname for the chatter via the response queue.
func (c *Chatter) getNickname() {
	c.sendResponse(ChatRspTypeGetNickname, c.Nickname(), nil)
}

// Nickname returns the raw nickname for the chatter.
func (c *Chatter) Nickname() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.nickname
}

// listRooms returns a list of chat rooms to the chatter.
func (c *Chatter) listRooms() {
	c.sendResponse(ChatRspTypeListRooms, "", c.cMngr.list())
}

// ChatterStats is a simple structure for returning statistic information on the chatter.
type ChatterStats struct {
	Nickname   string    `json:"nickname"`   // The nickname of the chatter.
	RemoteAddr string    `json:"remoteAddr"` // The remote IP and port of the chatter.
	Start      time.Time `json:"start"`      // The start time of the chatter.
	LastReq    time.Time `json:"lastReq"`    // The last request time from the chatter.
	LastRsp    time.Time `json:"lastRsp"`    // The last response time to the chatter.
	ReqCount   uint64    `json:"reqcount"`   // Total requests received.
	RspCount   uint64    `json:"rspCount"`   // Total responses sent.
}

// ChatterStatsNew returns status information on the chatter.
func (c *Chatter) ChatterStatsNew() *ChatterStats {
	c.mu.Lock()
	defer c.mu.Unlock()
	return &ChatterStats{
		Nickname:   c.nickname,
		RemoteAddr: c.ws.Request().RemoteAddr,
		Start:      c.start,
		LastReq:    c.lastReq,
		LastRsp:    c.lastRsp,
		ReqCount:   c.reqCount,
		RspCount:   c.rspCount,
	}
}

// sendRequestToRoom sends the request to a room or creates a mew room to receive the request.
func (c *Chatter) sendRequestToRoom(r *ChatRequest) {
	if r.RoomName == "" {
		c.sendResponse(ChatRspTypeErrRoomMandatory, "Room name is mandatory to access a room.", nil)
		return
	}
	m, err := c.cMngr.findCreate(r.RoomName)
	if err != nil {
		c.sendResponse(ChatRspTypeErrMaxRoomsReached, err.Error(), nil)
		return
	}
	m.reqq <- r
}

// sendResponse sends a message to the send() go routine to send message back to chatter.
func (c *Chatter) sendResponse(rt int, msg string, l []string) {
	if l == nil {
		l = []string{}
	}
	if rsp, err := ChatResponseNew("", rt, msg, l); err == nil {
		select {
		case <-c.done:
		default:
			c.rspq <- rsp
		}
	}
}
