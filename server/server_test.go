package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

const (
	testServerHostname   = "localhost"
	testServerPort       = 6660
	testServerMaxConns   = 4
	testServerMaxRooms   = 2
	testChatRoomName1    = "Room1"
	testChatRoomName2    = "Room2"
	testChatRoomName3    = "Room3"
	testChatterNickname1 = "ChatMonkey"
	testChatterNickname2 = "MonkeyTester"
)

var (
	testChatterStartTime   time.Time
	testChatterLastReqTime time.Time
	testChatterReqs        uint64
	testChatterRsps        uint64

	testRoomStartTime   time.Time
	testRoomLastReqTime time.Time
	testRoomReqs        uint64
	testRoomRsps        uint64

	testSrvr         *Server
	testSrvrURL      = fmt.Sprintf("ws://%s:%d/v1.0/chat", testServerHostname, testServerPort)
	testSrvrURLAlive = fmt.Sprintf("http://%s:%d/v1.0/alive", testServerHostname, testServerPort)
	testSrvrURLStats = fmt.Sprintf("http://%s:%d/v1.0/stats", testServerHostname, testServerPort)
	testSrvrOrg      = fmt.Sprintf("ws://%s/", testServerHostname)

	TestServerSetNickname = fmt.Sprintf(`{"reqType":%d,"content":"%s"}`,
		ChatReqTypeSetNickname, testChatterNickname1)
	TestServerSetNickname2 = fmt.Sprintf(`{"reqType":%d,"content":"%s"}`,
		ChatReqTypeSetNickname, testChatterNickname2)
	TestServerSetNicknameErr = fmt.Sprintf(`{"reqType":%d,"content":""}`,
		ChatReqTypeSetNickname)
	TestServerSetNicknameExp = fmt.Sprintf(`{"roomName":"","rspType":%d,"content":"Nickname set to \"%s\".","list":[]}`,
		ChatRspTypeSetNickname, testChatterNickname1)
	TestServerSetNicknameExpErr = fmt.Sprintf(`{"roomName":"","rspType":%d,"content":"Nickname cannot be blank.","list":[]}`,
		ChatRspTypeErrNicknameMandatory)

	TestServerGetNickname    = fmt.Sprintf(`{"reqType":%d}`, ChatReqTypeGetNickname)
	TestServerGetNicknameExp = fmt.Sprintf(`{"roomName":"","rspType":%d,"content":"%s","list":[]}`,
		ChatRspTypeGetNickname, testChatterNickname1)

	TestServerListRooms     = fmt.Sprintf(`{"reqType":%d}`, ChatReqTypeListRooms)
	TestServerListRoomsExp0 = fmt.Sprintf(`{"roomName":"","rspType":%d,"content":"","list":[]}`,
		ChatRspTypeListRooms)
	TestServerListRoomsExp1 = fmt.Sprintf(`{"roomName":"","rspType":%d,"content":"","list":["%s"]}`,
		ChatRspTypeListRooms, testChatRoomName1)

	TestServerJoin       = fmt.Sprintf(`{"roomName":"%s","reqType":%d}`, testChatRoomName1, ChatReqTypeJoin)
	TestServerJoinHidden = fmt.Sprintf(`{"roomName":"%s","reqType":%d,"content":"hidden"}`,
		testChatRoomName1, ChatReqTypeJoin)
	TestServerJoinErr = fmt.Sprintf(`{"roomName":"","reqType":%d}`, ChatReqTypeJoin)
	TestServerJoin2   = fmt.Sprintf(`{"roomName":"%s","reqType":%d}`, testChatRoomName2, ChatReqTypeJoin)
	TestServerJoin3   = fmt.Sprintf(`{"roomName":"%s","reqType":%d}`, testChatRoomName3, ChatReqTypeJoin)
	TestServerJoinExp = fmt.Sprintf(`{"roomName":"%s","rspType":%d,`+
		`"content":"%s has joined the room.","list":["%s"]}`, testChatRoomName1, ChatRspTypeJoin,
		testChatterNickname1, testChatterNickname1)
	TestServerJoinExp2 = fmt.Sprintf(`{"roomName":"%s","rspType":%d,`+
		`"content":"%s has joined the room.","list":["%s"]}`, testChatRoomName2, ChatRspTypeJoin,
		testChatterNickname1, testChatterNickname1)
	TestServerJoinExpHidden = fmt.Sprintf(`{"roomName":"%s","rspType":%d,`+
		`"content":"%s has joined the room.","list":[]}`, testChatRoomName1, ChatRspTypeJoin,
		testChatterNickname1)

	TestServerJoinExpErr = fmt.Sprintf(`{"roomName":"","rspType":%d,`+
		`"content":"Room name is mandatory to access a room.","list":[]}`, ChatRspTypeErrRoomMandatory)
	TestServerJoinExpErrX2 = fmt.Sprintf(`{"roomName":"%s","rspType":%d,`+
		`"content":"You are already a member of room \"%s\".","list":[]}`, testChatRoomName1,
		ChatRspTypeErrAlreadyJoined, testChatRoomName1)
	TestServerJoinExpErrSame = fmt.Sprintf(`{"roomName":"%s","rspType":%d,`+
		`"content":"Nickname \"%s\" is already in use in room \"%s\".","list":[]}`, testChatRoomName1,
		ChatRspTypeErrNicknameUsed, testChatterNickname1, testChatRoomName1)
	TestServerJoinExpErrRoom = fmt.Sprintf(`{"roomName":"","rspType":%d,`+
		`"content":"Maximum number of rooms reached. Cannot create new room.","list":[]}`,
		ChatRspTypeErrMaxRoomsReached)

	TestServerListNames     = fmt.Sprintf(`{"roomName":"%s","reqType":%d}`, testChatRoomName1, ChatReqTypeListNames)
	TestServerListNamesExp0 = fmt.Sprintf(`{"roomName":"%s","rspType":%d,"content":"","list":[]}`,
		testChatRoomName1, ChatRspTypeListNames)
	TestServerListNamesExp1 = fmt.Sprintf(`{"roomName":"%s","rspType":%d,"content":"","list":["%s"]}`,
		testChatRoomName1, ChatRspTypeListNames, testChatterNickname1)

	TestServerHideNickname    = fmt.Sprintf(`{"roomName":"%s","reqType":%d}`, testChatRoomName1, ChatReqTypeHide)
	TestServerHideNicknameExp = fmt.Sprintf(`{"roomName":"%s","rspType":%d,"content":"You are now hidden in room \"%s\".","list":[]}`,
		testChatRoomName1, ChatRspTypeHide, testChatRoomName1)
	TestServerUnhideNickname    = fmt.Sprintf(`{"roomName":"%s","reqType":%d}`, testChatRoomName1, ChatReqTypeUnhide)
	TestServerUnhideNicknameExp = fmt.Sprintf(`{"roomName":"%s","rspType":%d,"content":"You are now unhidden in room`+
		` \"%s\".","list":[]}`,
		testChatRoomName1, ChatRspTypeUnhide, testChatRoomName1)

	TestServerMsg = fmt.Sprintf(`{"roomName":"%s","reqType":%d,"content":"Hello you monkeys."}`,
		testChatRoomName1, ChatReqTypeMsg)
	TestServerMsgExp = fmt.Sprintf(`{"roomName":"%s","rspType":%d,"content":"%s: Hello you monkeys.","list":[]}`,
		testChatRoomName1, ChatRspTypeMsg, testChatterNickname1)
	TestServerMsgExpErrHide = fmt.Sprintf(`{"roomName":"%s","rspType":%d,"content":"Nickname \"%s\" `+
		`is hidden. Cannot post in room \"%s\".","list":[]}`,
		testChatRoomName1, ChatRspTypeErrHiddenNickname, testChatterNickname1, testChatRoomName1)

	TestServerLeave    = fmt.Sprintf(`{"roomName":"%s","reqType":%d}`, testChatRoomName1, ChatReqTypeLeave)
	TestServerLeaveExp = fmt.Sprintf(`{"roomName":"%s","rspType":%d,"content":"You have left room \"%s\".","list":[]}`,
		testChatRoomName1, ChatRspTypeLeave, testChatRoomName1)
)

func tTestIncrChatterStats() {
	testChatterLastReqTime = time.Now()
	testChatterReqs++
	testChatterRsps++
}

func tTestIncrRoomStats() {
	testRoomLastReqTime = time.Now()
	testRoomReqs++
	testRoomRsps++
}

func TestServerStartup(t *testing.T) {
	opts := &Options{
		Name:     "Test Server",
		Hostname: testServerHostname,
		Port:     testServerPort,
		ProfPort: 6060,
		MaxConns: testServerMaxConns,
		MaxRooms: testServerMaxRooms,
		MaxIdle:  0,
		MaxProcs: 1,
		Debug:    true,
	}
	runtime.GOMAXPROCS(1)
	testSrvr = New(opts)
	go func() { testSrvr.Start() }()
}

func TestServerPrintVersion(t *testing.T) {
	t.Parallel()
	t.Skip("Exit cannot be covered.")
}

func TestServerIsRunning(t *testing.T) {
	time.Sleep(2 * time.Second) // Make sure we are all ready.
	if !testSrvr.isRunning() {
		t.Errorf("Server should be runnning.")
	}
}

func TestServerValidWSSession(t *testing.T) {
	var rsp = make([]byte, 1024)
	var n int
	testChatterStartTime := time.Now()
	ws1, err := websocket.Dial(testSrvrURL, "", testSrvrOrg)
	if err != nil {
		t.Errorf("Server dialing error: %s", err)
		return
	}
	defer ws1.Close()

	ws2, err := websocket.Dial(testSrvrURL, "", testSrvrOrg)
	if err != nil {
		t.Errorf("Server dialing error: %s", err)
		return
	}
	defer ws2.Close()

	// Set Nickname
	tTestIncrChatterStats()
	if _, err := ws1.Write([]byte(TestServerSetNickname)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result := string(rsp[:n])
	if result != TestServerSetNicknameExp {
		t.Errorf("Set Nickname error.\nExpected: %s\n\nActual: %s\n", TestServerSetNicknameExp, result)
	}

	// Get Nickname
	tTestIncrChatterStats()
	if _, err := ws1.Write([]byte(TestServerGetNickname)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerGetNicknameExp {
		t.Errorf("Get Nickname error.\nExpected: %s\n\nActual: %s\n", TestServerGetNicknameExp, result)
	}

	// Get List of Rooms (0 rooms)
	tTestIncrChatterStats()
	if _, err := ws1.Write([]byte(TestServerListRooms)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerListRoomsExp0 {
		t.Errorf("Get list of rooms error.\nExpected: %s\n\nActual: %s\n", TestServerListRoomsExp0, result)
	}

	// Join a room
	tTestIncrChatterStats()
	testRoomStartTime := time.Now()
	tTestIncrRoomStats()
	if _, err := ws1.Write([]byte(TestServerJoin)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerJoinExp {
		t.Errorf("Join room error.\nExpected: %s\n\nActual: %s\n", TestServerJoinExp, result)
	}

	// Get List of Rooms (1 room)
	tTestIncrChatterStats()
	if _, err := ws1.Write([]byte(TestServerListRooms)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerListRoomsExp1 {
		t.Errorf("Get list of rooms error.\nExpected: %s\n\nActual: %s\n", TestServerListRoomsExp1, result)
	}

	// Get list of nicknames in a room (expect 1)
	tTestIncrChatterStats()
	tTestIncrRoomStats()
	if _, err := ws1.Write([]byte(TestServerListNames)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerListNamesExp1 {
		t.Errorf("Get list of names error.\nExpected: %s\n\nActual: %s\n", TestServerListNamesExp1, result)
	}

	// Hide nickname
	tTestIncrChatterStats()
	tTestIncrRoomStats()
	if _, err := ws1.Write([]byte(TestServerHideNickname)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerHideNicknameExp {
		t.Errorf("Hide nickname error.\nExpected: %s\n\nActual: %s\n", TestServerHideNicknameExp, result)
	}

	// Validate nickname is invisible in list (expect 0)
	tTestIncrChatterStats()
	tTestIncrRoomStats()
	if _, err := ws1.Write([]byte(TestServerListNames)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerListNamesExp0 {
		t.Errorf("Test hidden nickname. Get list of names error.\nExpected: %s\n\nActual: %s\n",
			TestServerListNamesExp0, result)
	}

	// Unhide nickname
	tTestIncrChatterStats()
	tTestIncrRoomStats()
	if _, err := ws1.Write([]byte(TestServerUnhideNickname)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerUnhideNicknameExp {
		t.Errorf("Unhide nickname error.\nExpected: %s\n\nActual: %s\n", TestServerUnhideNicknameExp, result)
	}

	// Validate nickname is now visible in list (expect 1)
	tTestIncrChatterStats()
	tTestIncrRoomStats()
	if _, err := ws1.Write([]byte(TestServerListNames)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerListNamesExp1 {
		t.Errorf("Test unhidden nickname. Get list of names error.\nExpected: %s\n\nActual: %s\n",
			TestServerListNamesExp1, result)
	}

	// Send a message to the room.
	tTestIncrChatterStats()
	tTestIncrRoomStats()
	if _, err := ws1.Write([]byte(TestServerMsg)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerMsgExp {
		t.Errorf("Send message error.\nExpected: %s\n\nActual: %s\n", TestServerMsgExp, result)
	}

	// Leave the room
	tTestIncrChatterStats()
	tTestIncrRoomStats()
	if _, err := ws1.Write([]byte(TestServerLeave)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerLeaveExp {
		t.Errorf("Leave room error.\nExpected: %s\n\nActual: %s\n", TestServerLeaveExp, result)
	}

	// Validate nickname is fully out of the room list (expect 0)
	tTestIncrChatterStats()
	testRoomLastReqTime := time.Now()
	tTestIncrRoomStats()
	if _, err := ws1.Write([]byte(TestServerListNames)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerListNamesExp0 {
		t.Errorf("Test leave room. Get list of names error.\nExpected: %s\n\nActual: %s\n",
			TestServerListNamesExp0, result)
	}

	// Join the room as hidden
	tTestIncrChatterStats()
	tTestIncrRoomStats()
	if _, err := ws1.Write([]byte(TestServerJoinHidden)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerJoinExpHidden {
		t.Errorf("Join room error.\nExpected: %s\n\nActual: %s\n", TestServerJoinExpHidden, result)
	}

	// Validate again nickname is invisible in list (expect 0)
	tTestIncrChatterStats()
	tTestIncrRoomStats()
	if _, err := ws1.Write([]byte(TestServerListNames)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerListNamesExp0 {
		t.Errorf("Test hidden nickname. Get list of names error.\nExpected: %s\n\nActual: %s\n",
			TestServerListNamesExp0, result)
	}

	// Leave the room
	tTestIncrChatterStats()
	tTestIncrRoomStats()
	if _, err := ws1.Write([]byte(TestServerLeave)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerLeaveExp {
		t.Errorf("Leave room error.\nExpected: %s\n\nActual: %s\n", TestServerLeaveExp, result)
	}

	// Join room again for stats test.
	tTestIncrChatterStats()
	tTestIncrRoomStats()
	if _, err := ws1.Write([]byte(TestServerJoin)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerJoinExp {
		t.Errorf("Join room error.\nExpected: %s\n\nActual: %s\n", TestServerJoinExp, result)
	}
	rm, _ := testSrvr.cMngr.find(testChatRoomName1)

	var ch *Chatter
	for k := range rm.chatters {
		ch = k
	}

	// Have chatter 2 also participate in the room for coverage.
	if _, err := ws2.Write([]byte(TestServerSetNickname2)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws2.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}

	tTestIncrRoomStats() // one for ws2
	testChatterRsps++    // but ws1 also gets a join update increment
	testRoomRsps++
	if _, err := ws2.Write([]byte(TestServerJoin)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws2.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}

	// Validate Chat Room statistics from this session.
	s := rm.ChatRoomStatsNew()
	if s.Name != testChatRoomName1 {
		t.Errorf("Room stats error. Name is incorrect. \nExpected: %s\n\nActual: %s\n",
			testChatRoomName1, s.Name)
	}
	if s.Start.Before(testRoomStartTime) || s.Start.Equal(testRoomStartTime) {
		t.Errorf("Room stats error. Start Time is out of range.")
	}
	if s.LastReq.Before(testRoomLastReqTime) || s.LastReq.Equal(testRoomLastReqTime) {
		t.Errorf("Room stats error. Last Request Time is out of range.")
	}
	if s.LastRsp.Before(testRoomLastReqTime) || s.LastRsp.Equal(testRoomLastReqTime) {
		t.Errorf("Room stats error. Last Response Time is out of range.")
	}
	if s.ReqCount != testRoomReqs {
		t.Errorf("Room stats error. ReqCount is incorrect.\nExpected: %d\n\nActual: %d\n",
			testRoomReqs, s.ReqCount)
	}
	if s.RspCount != testRoomRsps {
		t.Errorf("Room stats error. RsqCount is incorrect.\nExpected: %d\n\nActual: %d\n",
			testRoomRsps, s.RspCount)
	}

	// Validate Chatter statistics for this session.
	cs := ch.ChatterStatsNew()
	if cs.Nickname != testChatterNickname1 {
		t.Errorf("Chatter stats error. Nickname not correct.")
	}
	if cs.RemoteAddr == "" {
		t.Errorf("Chatter stats error. RemoteAddr not set.")
	}
	if cs.Start.Before(testChatterStartTime) || cs.Start.Equal(testChatterStartTime) {
		t.Errorf("Chatter stats error. Start Time is out of range.")
	}
	if cs.LastReq.Before(testChatterLastReqTime) || cs.LastReq.Equal(testChatterLastReqTime) {
		t.Errorf("Chatter stats error. Last Request Time is out of range.")
	}
	if cs.LastRsp.Before(testChatterLastReqTime) || cs.LastRsp.Equal(testChatterLastReqTime) {
		t.Errorf("Chatter stats error. Last Response Time is out of range.")
	}
	if cs.ReqCount != testChatterReqs {
		t.Errorf("Chatter stats error. ReqCount is incorrect.\nExpected: %d\n\nActual: %d\n",
			testChatterReqs, cs.ReqCount)
	}
	if cs.RspCount != testChatterRsps {
		t.Errorf("Chatter stats error. RsqCount is incorrect.\nExpected: %d\n\nActual: %d\n",
			testChatterRsps, cs.RspCount)
	}

	// Chatter 1 leaves again for coverage.
	tTestIncrChatterStats()
	tTestIncrRoomStats()
	if _, err := ws1.Write([]byte(TestServerLeave)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
}

func TestServerWSErrorSession(t *testing.T) {
	var rsp = make([]byte, 1024)
	var n int
	ws1, err := websocket.Dial(testSrvrURL, "", testSrvrOrg)
	if err != nil {
		t.Errorf("Server dialing error for ws1: %s", err)
		return
	}
	defer ws1.Close()
	ws2, err := websocket.Dial(testSrvrURL, "", testSrvrOrg)
	if err != nil {
		t.Errorf("Server dialing error for ws2: %s", err)
		return
	}
	defer ws2.Close()

	// Set nickname test error conditions.
	if _, err := ws1.Write([]byte(TestServerSetNicknameErr)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result := string(rsp[:n])
	if result != TestServerSetNicknameExpErr {
		t.Errorf("Set nickname did not receive an error.\nExpected: %s\n\nActual: %s\n",
			TestServerSetNicknameExpErr, result)
	}

	// Join a room test err conditions
	if _, err := ws1.Write([]byte(TestServerJoinErr)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerJoinExpErr {
		t.Errorf("Join room did not receive an error.\nExpected: %s\n\nActual: %s\n",
			TestServerJoinExpErr, result)
	}

	// Set nickname correctly for user 1
	if _, err := ws1.Write([]byte(TestServerSetNickname)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerSetNicknameExp {
		t.Errorf("Set nickname received an error.\nExpected: %s\n\nActual: %s\n",
			TestServerSetNicknameExp, result)
	}

	// Set nickname user 2 same as user 1
	if _, err := ws2.Write([]byte(TestServerSetNickname)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws2.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerSetNicknameExp {
		t.Errorf("Set nickname received an error.\nExpected: %s\n\nActual: %s\n",
			TestServerSetNicknameExp, result)
	}

	// User 1 joins room 1
	if _, err := ws1.Write([]byte(TestServerJoin)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerJoinExp {
		t.Errorf("Join received an error.\nExpected: %s\n\nActual: %s\n", TestServerJoinExp, result)
	}

	// User 1 tries to join room 1 again. Should not be allowed.
	if _, err := ws1.Write([]byte(TestServerJoin)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerJoinExpErrX2 {
		t.Errorf("Join 2X should have received an error.\nExpected: %s\n\nActual: %s\n",
			TestServerJoinExpErrX2, result)
	}

	// Hide user 1 nickname from room.
	if _, err := ws1.Write([]byte(TestServerHideNickname)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerHideNicknameExp {
		t.Errorf("Hide nickname received an error.\nExpected: %s\n\nActual: %s\n",
			TestServerHideNicknameExp, result)
	}

	// Posting ability should be disabled if name is hidden.
	if _, err := ws1.Write([]byte(TestServerMsg)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerMsgExpErrHide {
		t.Errorf("Message with hidded name should have received an error.\nExpected: %s\n\nActual: %s\n",
			TestServerMsgExpErrHide, result)
	}

	// Nickname already used in room should prevent joining. User 2 joins room 1 w/ same
	// name User 1
	if _, err := ws2.Write([]byte(TestServerJoin)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws2.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerJoinExpErrSame {
		t.Errorf("Join with same name should have received an error.\nExpected: %s\n\nActual: %s\n",
			TestServerJoinExpErrSame, result)
	}

	// Should not be able to grow room limitation.
	if _, err := ws1.Write([]byte(TestServerJoin2)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerJoinExp2 {
		t.Errorf("Join should not receive an error.\nExpected: %s\n\nActual: %s\n",
			TestServerJoinExp2, result)
	}

	if _, err := ws1.Write([]byte(TestServerJoin3)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	result = string(rsp[:n])
	if result != TestServerJoinExpErrRoom {
		t.Errorf("Join should have received an error.\nExpected: %s\n\nActual: %s\n",
			TestServerJoinExpErrRoom, result)
	}

	// Test Max timeout
	testSrvr.info.MaxIdle = 3 // Set to 3 seconds
	// Now, send a command to enable the new MaxIdle
	if _, err := ws1.Write([]byte(TestServerListRooms)); err != nil {
		if err != nil {
			t.Errorf("Websocket send error: %s", err)
			return
		}
	}
	if n, err = ws1.Read(rsp); err != nil {
		if err != nil {
			t.Errorf("Websocket receive error: %s", err)
			return
		}
	}
	// Sleep 5 seconds and try to send a command again.
	time.Sleep(5 * time.Second)
	if _, err := ws1.Write([]byte(TestServerListRooms)); err != nil {
		if err == nil {
			t.Errorf("Websocket should have been closed.")
		}
	}
	testSrvr.info.MaxIdle = 0 // Reset
}

func TestHTTPRoutes(t *testing.T) {
	client := &http.Client{}
	rq, _ := http.NewRequest("GET", testSrvrURLAlive, nil)
	rq.Header.Add("Accept", "application/json")
	rq.Header.Add("Accept-Encoding", "gzip, deflate")
	rq.Header.Add("Content-Type", "application/json")
	r, _ := client.Do(rq)
	b, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	body := string(b)
	if body != "" {
		t.Errorf("/alive body should be empty.")
	}
	if r.StatusCode != http.StatusOK {
		t.Errorf("/alive returned invalid status code.\nExpected: %d\n\nActual: %d\n",
			http.StatusOK, r.StatusCode)
	}

	rq, _ = http.NewRequest("GET", testSrvrURLStats, nil)
	rq.Header.Add("Accept", "application/json")
	rq.Header.Add("Accept-Encoding", "gzip, deflate")
	rq.Header.Add("Content-Type", "application/json")
	r, _ = client.Do(rq)
	b, _ = ioutil.ReadAll(r.Body)
	r.Body.Close()
	body = string(b)
	if body == "" {
		t.Errorf("/status body should not be empty.")
	}
	if r.StatusCode != http.StatusOK {
		t.Errorf("/status returned invalid status code.\nExpected: %d\n\nActual: %d\n",
			http.StatusOK, r.StatusCode)
	}
}

func TestServerTakeDown(t *testing.T) {
	time.Sleep(1 * time.Second) // allow all connections to leave cleanly from previous test.
	ws1, err := websocket.Dial(testSrvrURL, "", testSrvrOrg)
	if err != nil {
		t.Errorf("Server dialing error: %s", err)
		return
	}
	defer ws1.Close()
	time.Sleep(1 * time.Second) // Keep one up.
	testSrvr.Shutdown()
	testSrvr.Shutdown() // Coverage only
	if testSrvr.isRunning() {
		t.Errorf("Server should have shut down.")
	}
	testSrvr = nil
}
