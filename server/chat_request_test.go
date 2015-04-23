package server

import (
	"fmt"
	"testing"
)

var (
	testChatReqJSONResult = fmt.Sprintf(`{"roomName":"Room 237",`+
		`"reqType":%d,"content":"JonnyGoLucky"}`, ChatReqTypeSetNickname)
)

func TestChatReqNew(t *testing.T) {
	t.Parallel()
	_, err := ChatRequestNew(nil, "Room 237", ChatReqTypeSetNickname, "JonnyGoLucky")
	if err != nil {
		t.Errorf("Chat Request new should not have returned an error for valid low type.")
	}
	_, err = ChatRequestNew(nil, "Room 237", ChatReqTypeLeave, "JonnyGoLucky")
	if err != nil {
		t.Errorf("Chat Request new should not have returned an error for valid high type.")
	}
	_, err = ChatRequestNew(nil, "Room 237", ChatReqTypeSetNickname-1, "JonnyGoLucky")
	if err == nil {
		t.Errorf("Chat Request new should have returned an error for out of range low req type.")
	}

	_, err = ChatRequestNew(nil, "Room 237", ChatReqTypeLeave+1, "JonnyGoLucky")
	if err == nil {
		t.Errorf("Chat Request new should not have returned an error for out of range high req type.")
	}
}

func TestChatReqString(t *testing.T) {
	t.Parallel()
	r, _ := ChatRequestNew(nil, "Room 237", ChatReqTypeSetNickname, "JonnyGoLucky")
	actual := fmt.Sprint(r)
	if actual != testChatReqJSONResult {
		t.Errorf("Chat Request not converted to json string.\n\nExpected: %s\n\nActual: %s\n",
			testChatReqJSONResult, actual)
	}
}
