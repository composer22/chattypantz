var EventTypes = {
	CHANGE_EVENT: "change"
};

var ActionTypes = {
	LOGIN: "LOGIN",
	REFRESH_LOGIN: "REFRESH_LOGIN",
	SEND_MESSAGE: "SEND_MESSAGE",
	LOGOUT: "LOGOUT"
};

var Server = "ws://127.0.0.1:6660/v1.0/chat";

// Requests sent to the server.
var RequestTypes = {
  SET_NICKNAME: 101,
  GET_NICKNAME: 102,
  LIST_ROOMS: 103,
  JOIN: 104,
  LIST_NAMES: 105,
  HIDE: 106,
  UNHIDE: 107,
  MSG: 108,
  LEAVE: 109
};

// Responses coming from the server.
var ResponseTypes = {
  // Command responses
  SET_NICKNAME: 101,
  GET_NICKNAME: 102,
  LIST_ROOMS: 103,
  JOIN: 104,
  LIST_NAMES: 105,
  HIDE: 106,
  UNHIDE: 107,
  MSG: 108,
  LEAVE: 109,

  // Error Conditions
  ERR_ROOM_MANDATORY: 1001,
  ERR_MAX_ROOMS_REACHED: 1002,
  ERR_NICKNAME_MANDATORY: 1003,
  ERR_ALREADY_JOINED: 1004,
  ERR_NICKNAME_USED: 1005,
  ERR_HIDDEN_NICKNAME: 1006,
  ERR_NOT_IN_ROOM: 1007,
  ERR_UNKNOWN_REQUEST: 1008
};
