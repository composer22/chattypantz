package server

import (
	"fmt"
	"net"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

var (
	maxChatMsgs  = 1000                   // Maximum cap of msgs that can be queued to the msg queue.
	MaxChatSleep = 100 * time.Millisecond // How long to sleep between respq peeks.

)

// Chatter is a wrapper around a connection that represents one chat client on the server.
type Chatter struct {
	ws       *websocket.Conn    // The socket to the remote client.
	nickname string             // The friendly nickname to display in a conversation.
	rspq     chan *ChatResponse // A channel to receive information to send to the remote client.
	maxIdle  int                // Maximum idle time for a client
	mu       sync.Mutex         // For locking access to chatter attributes.
	log      *ChatLogger        // Application log for events.
	wg       sync.WaitGroup     // wait group for all thread adds
}

// ChatterNew is a factory function that returns a new Chatter instance
func ChatterNew(c *websocket.Conn, mi int, cl *ChatLogger, addedOpts ...func(*Chatter)) *Chatter {
	ctr := &Chatter{
		ws:      c,
		rspq:    make(chan *ChatResponse, maxChatMsgs),
		maxIdle: mi,
		log:     cl,
	}
	// Additional hook for specialized custom options.
	for _, f := range addedOpts {
		f(ctr)
	}
	return ctr
}

// Run starts the event loop that manages the sending and receiving of information to the client.
func (c *Chatter) Run() {
	c.wg.Add(1)
	go c.send() // Spawn this in the background to send info to the client.
	c.receive() // Then wait on incoming commands.
	c.wg.Wait()
}

// receive polls and handles any commands or information sent from the remote client.

func (c *Chatter) receive() {
	var cmd string
	remoteAddr := fmt.Sprint(c.ws.Request().RemoteAddr)
	for {
		// Set optional idle timeout.
		if c.maxIdle > 0 {
			c.ws.SetReadDeadline(time.Now().Add(time.Duration(c.maxIdle) * time.Second))
		}

		if err := websocket.Message.Receive(c.ws, &cmd); err != nil {
			e, ok := err.(net.Error)
			switch {
			case ok && e.Timeout():
				c.log.LogSession("disconnected", remoteAddr, "Client forced to disconnect due to inactivity.")
			case err.Error() == "EOF":
				c.log.LogSession("disconnected", remoteAddr, "Client disconnected.")
			default:
				c.log.LogError(remoteAddr, fmt.Sprintf("Couldn't receive. Error: %s", err.Error()))
			}
			close(c.rspq) // signal to send() to stop.
			return
		}
		c.log.LogSession("received", remoteAddr, cmd)
		// unmarshal
		// send an error to client if cannot trap msg

		// TODO do something with the cmd
		// change nickname
		// get nickname
		// get room info: list of rooms, number of chatters, your current room
		// list members in room
		// join room
		// send message to room
		// leave room
		//

		time.Sleep(chatterMaxSleep) // Sleep before peeking again.
	}
}

// send is a go routine used to poll queued messages to send information to the client.
func (c *Chatter) send() {
	defer c.wg.Done()
	remoteAddr := fmt.Sprint(c.ws.Request().RemoteAddr)
loop:
	for {
		select {
		case rs, ok := <-c.rspq:
			if !ok { // Assume ch closed and shutdown notification
				c.ws.Close() // we need to close the channel to break the receive looper.
				break loop
			}
			// Send the packet. Assume it's JSON already.
			if err := websocket.Message.Send(c.ws, rs); err != nil {
				switch {
				case err.Error() == "EOF":
					c.log.LogSession("disconnected", remoteAddr, "Client disconnected.")
					break loop
				default:
					c.log.LogError(remoteAddr, fmt.Sprintf("Couldn't send. Error: %s", err.Error()))
				}
			}
		default:
			time.Sleep(MaxChatSleep) // Sleep before peeking again.
		}
	}
}
