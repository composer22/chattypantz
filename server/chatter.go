package server

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

var (
	maxChatterRsp   = 1000                   // The maximum number of responses in the response channel.
	maxChatterSleep = 100 * time.Millisecond // How long to sleep between chan peeks.
)

// Chatter is a wrapper around a connection that represents one chat client on the server.
type Chatter struct {
	srvr     *Server            // The server this chatter is connected to.
	ws       *websocket.Conn    // The socket to the remote client.
	nickname string             // The friendly nickname to display in a conversation.
	start    time.Time          // The start time of the connection.
	lastReq  time.Time          // The last request time of the connection.
	lastRsp  time.Time          // The last response time to the connection.
	reqCount int64              // Total requests received.
	rspCount int64              // Total responses sent.
	rspq     chan *ChatResponse // A channel to receive information to send to the remote client.
	mu       sync.Mutex         // For locking access to chatter attributes.
	wg       sync.WaitGroup     // Synchronization of channel close.
}

// ChatterNew is a factory function that returns a new Chatter instance
func ChatterNew(s *Server, c *websocket.Conn) *Chatter {
	return &Chatter{
		srvr: s,
		ws:   c,
		rspq: make(chan *ChatResponse, maxChatterRsp),
	}
}

// Run starts the event loop that manages the sending and receiving of information to the client.
func (c *Chatter) Run() {
	c.start = time.Now()
	c.srvr.wg.Add(1) // We let the big boss also perform waits for chatters,
	c.wg.Add(1)      //   but we also have our own.
	go c.send()      // Spawn responses in the background to the client.
	c.receive()      // Then wait on incoming requests.
}

// receive polls and handles any commands or information sent from the remote client.
func (c *Chatter) receive() {
	remoteAddr := fmt.Sprint(c.ws.Request().RemoteAddr)
	for {
		// Set optional idle timeout.
		if c.srvr.info.MaxIdle > 0 {
			c.ws.SetReadDeadline(time.Now().Add(time.Duration(c.srvr.info.MaxIdle) * time.Second))
		}
		var req ChatRequest
		if err := websocket.JSON.Receive(c.ws, &req); err != nil {
			e, ok := err.(net.Error)
			switch {
			case ok && e.Timeout():
				c.srvr.log.LogSession("disconnected", remoteAddr, "Client forced to disconnect due to inactivity.")
				c.shutDown()
			case err.Error() == "EOF":
				c.srvr.log.LogSession("disconnected", remoteAddr, "Client disconnected.")
				c.shutDown()
			case strings.Contains(err.Error(), "use of closed network connection"):
				c.shutDown()
			default:
				c.srvr.log.LogError(remoteAddr, fmt.Sprintf("Couldn't receive. Error: %s", err.Error()))
			}
			return
		}
		c.lastReq = time.Now()
		c.reqCount++
		c.srvr.log.LogSession("received", remoteAddr, fmt.Sprintf("%s", &req))
		switch req.ReqType {
		case ChatReqTypeSetNickname:
			c.setNickname(&req)
		case ChatReqTypeGetNickname:
			c.getNickname()
		case ChatReqTypeListRooms:
			c.listRooms()
		default: // let room handle the request or give error if no room name provided.
			req.Who = c
			c.sendRequestToRoom(&req)
		}
	}
}

// send is a go routine used to poll queued messages to send information to the client.
func (c *Chatter) send() {
	defer c.srvr.wg.Done()
	defer c.wg.Done()

	remoteAddr := fmt.Sprint(c.ws.Request().RemoteAddr)
	for {
		select {
		case <-c.srvr.done: // Server shutdown signal.
			c.ws.Close() // Break the receiv() loop and force a chatter shutdown.
			return
		case rsp, ok := <-c.rspq:
			if !ok { // Assume ch closed is a shutdown notification from anybody.
				c.ws.Close() // Break the receive() looper and force a chatter shutdown.
				return
			}
			c.lastRsp = time.Now()
			c.rspCount++
			c.srvr.log.LogSession("sent", remoteAddr, fmt.Sprintf("%s", rsp))
			if err := websocket.JSON.Send(c.ws, rsp); err != nil {
				switch {
				case err.Error() == "EOF":
					c.srvr.log.LogSession("disconnected", remoteAddr, "Client disconnected.")
					return
				default:
					c.srvr.log.LogError(remoteAddr, fmt.Sprintf("Couldn't send. Error: %s", err.Error()))
				}
			}
		default:
			time.Sleep(maxChatterSleep)
		}
	}
}

// shutDown removes the chatter from any rooms, and shuts down sending/receiving.
func (c *Chatter) shutDown() {
	// Clear out of all rooms and stop response signalling.
	c.srvr.roomMngr.removeChatterAllRooms(c)
	c.mu.Lock()
	close(c.rspq) // Signal to send() to stop.
	c.mu.Unlock()
	c.wg.Wait()
}

// setNickname sets the nickname for the chatter.
func (c *Chatter) setNickname(r *ChatRequest) {
	if r.Content == "" {
		c.sendResponse(ChatRspTypeErrNicknameMandatory, "Nickname cannot be blank.")
		return
	}
	c.nickname = r.Content
	c.sendResponse(ChatRspTypeSetNickname, fmt.Sprintf(`Nickname set to "%s".`, c.nickname))
}

// nickname returns the nickname for the chatter.
func (c *Chatter) getNickname() {
	c.sendResponse(ChatRspTypeNickname, c.nickname)
}

// listRooms returns a list of chat rooms to the chatter.
func (c *Chatter) listRooms() {
	c.sendResponse(ChatRspTypeListRooms, c.srvr.roomMngr.list())
}

// sendRequestToRoom sends the request to a room or creates a mew room to receive the request.
func (c *Chatter) sendRequestToRoom(r *ChatRequest) {
	if r.RoomName == "" {
		c.sendResponse(ChatRspTypeErrRoomIsMandatory, "Room name is mandatory to join a room.")
		return
	}
	m, err := c.srvr.roomMngr.createOrFind(r.RoomName)
	if err != nil {
		c.sendResponse(ChatRspTypeErrMaxRoomsReached, err.Error())
		return
	}
	m.reqq <- r
}

// sendResponse sends a message to a chatter.
func (c *Chatter) sendResponse(rt int, msg string) {
	c.mu.Lock()
	if c.rspq != nil {
		c.rspq <- ChatResponseNew("", rt, msg)
	}
	c.mu.Unlock()
}
