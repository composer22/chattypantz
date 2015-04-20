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
func ChatRoomManagerNew(x int, cl *ChatLogger) *ChatRoomManager {
	return &ChatRoomManager{
		rooms:    make(map[string]*ChatRoom),
		maxRooms: x,
		log:      cl,
	}
}

// list returns a list of chat room names.
func (m *ChatRoomManager) list() string {
	var names []string
	for n := range m.rooms {
		names = append(names, n)
	}
	return fmt.Sprint(names)
}

// createOrFind returns a new chat room for a given name or returns one already created.
func (m *ChatRoomManager) createOrFind(n string) (*ChatRoom, error) {
	r, ok := m.rooms[n]
	if !ok {
		if m.maxRooms > 0 && m.maxRooms == len(m.rooms) {
			return nil, errors.New("Maximum number of rooms reached. Cannot create new room.")
		}
		r = ChatRoomNew(n, m.log, &m.wg)
		m.wg.Add(1)
		go r.Run()
	}
	return r, nil
}

// shutDownRooms releases all rooms from processing and memory.
func (m *ChatRoomManager) shutDownRooms() {
	// Close the channel which signals a stop run
	for _, cr := range m.rooms {
		close(cr.reqq)
	}
	m.wg.Wait()
	m.rooms = nil
	m.rooms = make(map[string]*ChatRoom)
}

// removeChatterAllRooms releases the chatter from any rooms.
func (m *ChatRoomManager) removeChatterAllRooms(c *Chatter) {
	// Close the channel which signals a stop run
	for _, cr := range m.rooms {
		cr.reqq <- ChatRequestNew(c, cr.name, ChatReqTypeLeave, "")
	}
}
