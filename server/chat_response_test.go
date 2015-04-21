package server

import (
	"fmt"
	"testing"
)

var (
	testChatRspJSONResult = fmt.Sprintf(`{"roomName":"Room 237","rspType":%d,`+
		`"content":"JonnyGoLucky","list":["One","Two"]}`, ChatRspTypeSetNickname)
)

func TestChatRspNew(t *testing.T) {
	t.Parallel()

	_, err := ChatResponseNew("Room 237", ChatRspTypeSetNickname, "JonnyGoLucky", []string{"One", "Two"})
	if err != nil {
		t.Errorf("Chat Response new should not have returned an error for valid low type.")
	}
	_, err = ChatResponseNew("Room 237", ChatRspTypeLeave, "JonnyGoLucky", []string{"One", "Two"})
	if err != nil {
		t.Errorf("Chat Request new should not have returned an error for valid high type.")
	}
	_, err = ChatResponseNew("Room 237", ChatRspTypeErrRoomMandatory, "JonnyGoLucky", []string{"One", "Two"})
	if err != nil {
		t.Errorf("Chat Response new should not have returned an error for valid low err type.")
	}
	_, err = ChatResponseNew("Room 237", ChatRspTypeErrUnknownReq, "JonnyGoLucky", []string{"One", "Two"})
	if err != nil {
		t.Errorf("Chat Request new should have returned an error for valid high err type.")
	}

	_, err = ChatResponseNew("Room 237", ChatRspTypeSetNickname-1, "JonnyGoLucky", []string{"One", "Two"})
	if err == nil {
		t.Errorf("Chat Response new should have returned an error for invalid low type.")
	}
	_, err = ChatResponseNew("Room 237", ChatRspTypeLeave+1, "JonnyGoLucky", []string{"One", "Two"})
	if err == nil {
		t.Errorf("Chat Request new should have returned an error for invalid high type.")
	}
	_, err = ChatResponseNew("Room 237", ChatRspTypeErrRoomMandatory-1, "JonnyGoLucky", []string{"One", "Two"})
	if err == nil {
		t.Errorf("Chat Response new should have returned an error for invalid low err type.")
	}
	_, err = ChatResponseNew("Room 237", ChatRspTypeErrUnknownReq+1, "JonnyGoLucky", []string{"One", "Two"})
	if err == nil {
		t.Errorf("Chat Request new should have returned an error for invalid high err type.")
	}
}

func TestChatRspString(t *testing.T) {
	t.Parallel()
	r, _ := ChatResponseNew("Room 237", ChatRspTypeSetNickname, "JonnyGoLucky", []string{"One", "Two"})
	actual := fmt.Sprint(r)
	if actual != testChatRspJSONResult {
		t.Errorf("Chat Response not converted to json string.\n\nExpected: %s\n\nActual: %s\n",
			testChatRspJSONResult, actual)
	}
}
