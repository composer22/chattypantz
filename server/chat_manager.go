package server

import (
	"errors"
	"sync"

	"golang.org/x/net/websocket"
)

var (
	chatManagerErrMaxRooms     = errors.New("maximum number of rooms reached")
	chatManagerErrRoomExists   = errors.New("room already exists")
	chatManagerErrRoomNotEmpty = errors.New("room is not empty")
	chatManagerErrRoomNotFound = errors.New("chatroom not found")
)

// ChatManager represents a control hub of chat rooms and chatters for the server.
type ChatManager struct {
	mu       sync.RWMutex         // Lock for update.
	rooms    map[string]*ChatRoom // A list of rooms on the server.
	chatters map[*Chatter]bool    // A list of chatters on the server.
	maxRooms int                  // Maximum number of rooms allowed to be created.
	maxIdle  int                  // Maximum idle time allowed for a ws connection.

	done chan bool      // Shut down chatters and rooms
	log  *ChatLogger    // Application log for events.
	wg   sync.WaitGroup // Synchronizer for manager reqq.
}

// ChatManagerNew is a factory function that returns a new instance of a chat manager.
func ChatManagerNew(maxr int, maxi int, l *ChatLogger) *ChatManager {
	return &ChatManager{
		rooms:    make(map[string]*ChatRoom),
		chatters: make(map[*Chatter]bool),
		maxRooms: maxr,
		maxIdle:  maxi,
		done:     make(chan bool),
		log:      l,
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
func (m *ChatManager) find(name string) (*ChatRoom, error) {
	m.mu.RLock()
	rm, ok := m.rooms[name]
	m.mu.RUnlock()
	if !ok {
		return nil, chatManagerErrRoomNotFound
	}
	return rm, nil
}

// findCreate returns a chat room for a given name or create a new one.
func (m *ChatManager) findCreate(name string) (*ChatRoom, error) {
	room, err := m.find(name)
	if err == nil {
		return room, nil
	}
	return m.createRoom(name)
}

// createRoom returns a new chat room,
func (m *ChatManager) createRoom(name string) (*ChatRoom, error) {
	_, err := m.find(name)
	if err == nil {
		return nil, chatManagerErrRoomExists
	}
	maxr := m.MaxRooms()
	m.mu.Lock() // cover rooms
	if maxr > 0 && maxr == len(m.rooms) {
		m.mu.Unlock()
		return nil, chatManagerErrMaxRooms
	}
	room := ChatRoomNew(name, m.done, m.log, &m.wg)
	m.rooms[name] = room
	m.wg.Add(1)
	go room.Run()
	m.mu.Unlock()
	return room, nil
}

// renameRoom is used to change the name of a room.
func (m *ChatManager) renameRoom(oldName string, newName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	room, ok := m.rooms[oldName]
	if !ok {
		return chatManagerErrRoomNotFound
	}
	_, ok = m.rooms[newName]
	if ok {
		return chatManagerErrRoomExists
	}
	if !room.isEmpty() {
		return chatManagerErrRoomNotEmpty
	}

	delete(m.rooms, oldName)
	room.SetName(newName)
	m.rooms[newName] = room
	return nil
}

// deleteRoom stops a chat room from running and removes it from the directory.
func (m *ChatManager) deleteRoom(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	room, ok := m.rooms[name]
	if !ok {
		return chatManagerErrRoomNotFound
	}
	if !room.isEmpty() {
		return chatManagerErrRoomNotEmpty
	}
	delete(m.rooms, name)
	close(room.reqq)
	return nil
}

// removeChatterAllRooms sends a broadcast to all rooms to release the chatter.
func (m *ChatManager) removeChatterAllRooms(c *Chatter) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, r := range m.rooms {
		if req, err := ChatRequestNew(c, r.Name(), ChatReqTypeLeave, ""); err == nil {
			r.reqq <- req
		}
	}
}

// getRoomStats returns statistics from each room.
func (m *ChatManager) getRoomStats() []*ChatRoomStats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var stats = []*ChatRoomStats{}
	for _, r := range m.rooms {
		stats = append(stats, r.ChatRoomStatsNew())
	}
	return stats
}

// registerChatter registers a new chatter with the chat manager.
func (m *ChatManager) registerNewChatter(ws *websocket.Conn) *Chatter {
	m.mu.Lock()
	defer m.mu.Unlock()
	chatr := ChatterNew(m, ws, m.log)
	m.chatters[chatr] = true
	return chatr
}

// getChatterStats returns statistics from all chatters
func (m *ChatManager) getChatterStats() []*ChatterStats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var stats = []*ChatterStats{}
	for c := range m.chatters {
		stats = append(stats, c.ChatterStatsNew())
	}
	return stats
}

// unregisterChatter removes a new chatter from the chat manager.
func (m *ChatManager) unregisterChatter(c *Chatter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.chatters[c]; ok {
		delete(m.chatters, c)
	}
}

// Shuts down the chatters and the rooms. Used by server on quit.
func (m *ChatManager) shutdownAll() {
	close(m.done)
	m.wg.Wait()
	m.mu.Lock()
	m.rooms = make(map[string]*ChatRoom)
	m.chatters = make(map[*Chatter]bool)
	m.mu.Unlock()
}

// MaxRooms returns the current maximum number of rooms allowed on the server.
func (m *ChatManager) MaxRooms() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.maxRooms
}

// SetMaxRooms sets the maximum number of rooms allowed on the server.
func (m *ChatManager) SetMaxRooms(maxr int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.maxRooms = maxr
}

// MaxIdle returns the current maximum idle time for a connection.
func (m *ChatManager) MaxIdle() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.maxIdle
}

// SetMaxIdle sets the maximum idle time for a connection.
func (m *ChatManager) SetMaxIdle(maxi int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.maxIdle = maxi
}
