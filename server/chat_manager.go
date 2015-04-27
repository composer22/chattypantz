package server

import (
	"errors"
	"fmt"
	"sync"

	"golang.org/x/net/websocket"
)

// ChatManager represents a control hub of chat rooms and chatters for the server.
type ChatManager struct {
	rooms    map[string]*ChatRoom // A list of rooms on the server.
	chatters map[*Chatter]bool    // A list of chatters on the server.
	done     chan bool            // Shut down chatters and rooms
	maxRooms int                  // Maximum number of rooms allowed to be created.
	maxIdle  int                  // Maximum idle time allowed for a ws connection.
	log      *ChatLogger          // Application log for events.
	wg       sync.WaitGroup       // Synchronizer for manager reqq.
	mu       sync.Mutex           // Lock for update.
}

// ChatManagerNew is a factory function that returns a new instance of a chat manager.
func ChatManagerNew(n int, mi int, cl *ChatLogger) *ChatManager {
	return &ChatManager{
		rooms:    make(map[string]*ChatRoom),
		chatters: make(map[*Chatter]bool),
		done:     make(chan bool),
		maxRooms: n,
		maxIdle:  mi,
		log:      cl,
	}
}

// list returns a list of chat room names.
func (m *ChatManager) list() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	var names []string
	for n := range m.rooms {
		names = append(names, n)
	}
	return names
}

// find will find a chat room for a given name.
func (m *ChatManager) find(n string) (*ChatRoom, error) {
	m.mu.Lock()
	r, ok := m.rooms[n]
	m.mu.Unlock()
	if !ok {
		return nil, errors.New(fmt.Sprintf(`Chatroom "%s" not found.`, n))
	}
	return r, nil
}

// findCreate returns a chat room for a given name or create a new one.
func (m *ChatManager) findCreate(n string) (*ChatRoom, error) {
	r, err := m.find(n)
	if err != nil {
		m.mu.Lock() // cover rooms
		if m.maxRooms > 0 && m.maxRooms == len(m.rooms) {
			m.mu.Unlock()
			return nil, errors.New("Maximum number of rooms reached. Cannot create new room.")
		}
		r = ChatRoomNew(n, m.done, m.log, &m.wg)
		m.rooms[n] = r
		m.wg.Add(1)
		go r.Run()
		m.mu.Unlock()
	}
	return r, nil
}

// removeChatterAllRooms releases the chatter from any rooms.
func (m *ChatManager) removeChatterAllRooms(c *Chatter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, r := range m.rooms {
		if q, err := ChatRequestNew(c, r.name, ChatReqTypeLeave, ""); err == nil {
			r.reqq <- q
		}
	}
}

// getRoomStats returns statistics from each room.
func (m *ChatManager) getRoomStats() []*ChatRoomStats {
	m.mu.Lock()
	defer m.mu.Unlock()
	var s = []*ChatRoomStats{}
	for _, r := range m.rooms {
		s = append(s, r.ChatRoomStatsNew())
	}
	return s
}

// registerChatter registers a new chatter with the chat manager.
func (m *ChatManager) registerNewChatter(ws *websocket.Conn) *Chatter {
	m.mu.Lock()
	defer m.mu.Unlock()
	c := ChatterNew(m, ws, m.maxIdle, m.log)
	m.chatters[c] = true
	return c
}

// getChatterStats returns statistics from all chatters
func (m *ChatManager) getChatterStats() []*ChatterStats {
	m.mu.Lock()
	defer m.mu.Unlock()
	var s = []*ChatterStats{}
	for c := range m.chatters {
		s = append(s, c.ChatterStatsNew())
	}
	return s
}

// unregisterChatter removes a new chatter from the chat manager.
func (m *ChatManager) unregisterChatter(c *Chatter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.chatters[c]; ok {
		delete(m.chatters, c)
	}
}

// Shuts down the chatters and the rooms.
func (m *ChatManager) shutdownAll() {
	close(m.done)
	m.wg.Wait()
	m.rooms = make(map[string]*ChatRoom)
	m.chatters = make(map[*Chatter]bool)
}
