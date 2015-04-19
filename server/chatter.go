package server

import (
	"fmt"
	"net"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

var (
	maxChatMsgs     = 200                    // Maximum cap of msgs that can be queued to the msg queue.
	chatterMaxSleep = 250 * time.Millisecond // How long to sleep between msgq peeks.

)

// Chatter is a wrapper around a connection that represents one chat client on the server.
type Chatter struct {
	Conn        *net.Conn      // The socket to the remote client.
	Nickname    string         // The friendly nickname to display in a conversation.
	CurrentRoom *ChatRoom      // The current chatroom they are in.
	Msgq        chan string    // A channel to receive information to send to the remote client.
	maxIdle     int            // Maximum idle time for a client
	mu          sync.Mutex     // For locking access to chatter attributes.
	log         *ChatLogger    // Application log for events.
	wg          sync.WaitGroup // wait group for all thread adds
}

// ChatterNew is a factory function that returns a new Chatter instance
func ChatterNew(c *net.Conn, mi int, cl *ChatLogger, addedOpts ...func(*Chatter)) *Chatter {
	ctr := &Chatter{
		Conn:    c,
		Msgq:    make(chan string, maxChatMsgs),
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
		if c.MaxIdle > 0 {
			ws.SetReadDeadline(time.Now().Add(time.Duration(c.MaxIdle) * time.Second))
		}

		if err := websocket.Message.Receive(ws, &cmd); err != nil {
			e, ok := err.(net.Error)
			switch {
			case ok && e.Timeout():
				c.log.LogSession("disconnected", remoteAddr, "Client forced to disconnect due to inactivity.")
			case err.Error() == "EOF":
				c.log.LogSession("disconnected", remoteAddr, "Client disconnected.")
			default:
				c.log.LogError(remoteAddr, fmt.Sprintf("Couldn't receive. Error: %s", err.Error()))
			}
			close(c.Msgq) // signal to send() to stop.
			return
		}
		c.log.LogSession("received", remoteAddr, cmd)
		// unmarshal
		// send an error to client if cannot trap msg

		// TODO do something with the cmd
		// change nickname
		// get nickname
		// get room info: list of rooms, number of chatters, your current room
		// enter room (no if already there)
		// leave room (no if not in room)
		// send message to room  (no if not in room)
		time.Sleep(chatterMaxSleep) // Sleep before peeking again.
	}
}

// send is a go routine used to poll queued messages to send information to the client.
func (c *Chatter) send() {
	defer c.wg.Done()
loop:
	for {
		select {
		case msg, ok := <-c.Msgq:
			if !ok { // Assume ch closed and shutdown notification
				c.ws.Close() // we need to close the channel to break the receive looper.
				break loop
			}
			// Send the packet. Assume it's JSON already.
			if err := websocket.Message.Send(c.ws, msg); err != nil {
				switch {
				case err.Error() == "EOF":
					c.log.LogSession("disconnected", remoteAddr, "Client disconnected.")
					break loop
				default:
					c.log.LogError(remoteAddr, fmt.Sprintf("Couldn't send. Error: %s", err.Error()))
				}
			}
		default:
			time.Sleep(chatterMaxSleep) // Sleep before peeking again.
		}
	}
}
