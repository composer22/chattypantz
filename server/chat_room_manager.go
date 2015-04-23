package server

import (
	"errors"
	"fmt"
	"sync"
)

// ChatRoomManager represents a hub of chat rooms for the server.
type ChatRoomManager struct {
	rooms    map[string]*ChatRoom // A list of rooms on the server.
	maxRooms int                  // Maximum number of rooms allowed to be created
	log      *ChatLogger          // Application log for events.
	wg       sync.WaitGroup       // Synchronizer for manager reqq
}

// ChatRoomManagerNew is a factory function that returns a new instance of a chat room manager.
func ChatRoomManagerNew(n int, cl *ChatLogger) *ChatRoomManager {
	return &ChatRoomManager{
		rooms:    make(map[string]*ChatRoom),
		maxRooms: n,
		log:      cl,
	}
}

// list returns a list of chat room names.
func (m *ChatRoomManager) list() []string {
	var names []string
	for n := range m.rooms {
		names = append(names, n)
	}
	return names
}

// find will find a chat room for a given name.
func (m *ChatRoomManager) find(n string) (*ChatRoom, error) {
	r, ok := m.rooms[n]
	if !ok {
		return nil, errors.New(fmt.Sprintf(`Chatroom "%s" not found.`, n))
	}
	return r, nil
}

// findCreate returns a chat room for a given name or create a new one.
func (m *ChatRoomManager) findCreate(n string) (*ChatRoom, error) {
	r, err := m.find(n)
	if err != nil {
		if m.maxRooms > 0 && m.maxRooms == len(m.rooms) {
			return nil, errors.New("Maximum number of rooms reached. Cannot create new room.")
		}
		r = ChatRoomNew(n, m.log, &m.wg)
		m.rooms[n] = r
		m.wg.Add(1)
		go r.Run()
	}
	return r, nil
}

// removeChatterAllRooms releases the chatter from any rooms.
func (m *ChatRoomManager) removeChatterAllRooms(c *Chatter) {
	for _, r := range m.rooms {
		if q, err := ChatRequestNew(c, r.name, ChatReqTypeLeave, ""); err == nil {
			r.reqq <- q
		}
	}
}

// shutDownRooms releases all rooms from processing and memory.
func (m *ChatRoomManager) shutDownRooms() {
	// Close the channel which signals a stop run
	for _, r := range m.rooms {
		close(r.reqq)
	}
	m.wg.Wait()
	m.rooms = nil
	m.rooms = make(map[string]*ChatRoom)
}
