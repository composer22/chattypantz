// Tracks information about the connection to the server and the display

var ConnectionStore = Object.assign({}, EventEmitter.prototype, {
   chat: {
    socket: null,
    status: {
      error: '',
    },
    data: {
      users: [],
      room: 'Demo',
      nickname: '',
      history: '',
    }
  },

  emitChange: function() {
    this.emit(EventTypes.CHANGE_EVENT);
  },

  addChangeListener: function(callback) {
    this.on(EventTypes.CHANGE_EVENT, callback);
  },

  removeChangeListener: function(callback) {
    this.removeListener(EventTypes.CHANGE_EVENT, callback);
  },

  login: function() {
	 this.chat.init();
  },

  setStatusError: function(error) {
	 this.chat.status.error = error;
  },

  getStatusError: function() {
	 return this.chat.status.error;
  },

  setNickname: function(nickname) {
	 this.chat.data.nickname = nickname;
  },

  getNickname: function() {
	 return this.chat.data.nickname;
  }

});

ConnectionStore.chat.init = function() {
	this.data.users = [];
	this.data.history = '';
    // Open a socket and register event handlers.
    this.socket = new WebSocket(Server);
    this.socket.onopen = this.onOpenWS;
    this.socket.onmessage = this.onMessageWS;
    this.socket.onerror = this.onErrorWS;
    this.socket.onclose = this.onErrorWS;
};

//////// SOCKET EVENT HANDLERS ////////

ConnectionStore.chat.onOpenWS = function(e) {
	var cs = ConnectionStore;
	var csc = cs.chat;
    cs.setNickname(LoginStore.getNickname());
	cs.setStatusError('');
    csc.sendRequest('', RequestTypes.SET_NICKNAME,  csc.data.nickname);
    csc.sendRequest(csc.data.room, RequestTypes.JOIN, csc.data.room);
	React.render(
	 <ConnectedSection />,
	  document.getElementById("chattypantzapp")
	);
};

// Response arrived from the server.
ConnectionStore.chat.onMessageWS = function(message) {
  var csc = ConnectionStore.chat;
  var response = JSON.parse(message.data);
  switch (response.rspType) {
    case ResponseTypes.SET_NICKNAME:
      csc.data.history += "Chattypantz server: " + response.content + '\n';
	  ConnectionActions.refresh();
      break;
    case ResponseTypes.JOIN:
      csc.data.users = response.list;
      csc.data.history += "Chattypantz server: " + response.content + '\n';
	  ConnectionActions.refresh();
      break;
    case ResponseTypes.LIST_NAMES:
      csc.data.users = response.list;
	  ConnectionActions.refresh();
      break;
    case ResponseTypes.MSG:
      csc.data.history += response.content + '\n';
	  ConnectionActions.refresh();
      break;
    case ResponseTypes.LEAVE:
      csc.data.users = response.list;
      csc.data.history += "Chattypantz server: " + response.content + '\n';
	  ConnectionActions.refresh();
      break;
    case ResponseTypes.ERR_NICKNAME_USED:
      ConnectionStore.setStatusError(response.content);
      ConnectionActions.logout();
      break;
    default:
	  ConnectionStore.setStatusError(response.content);
	  ConnectionActions.logout();
  }
};

// Connection error.
ConnectionStore.chat.onErrorWS = function(e) {
  var nickname = ConnectionStore.getNickname();
  var err = ConnectionStore.getStatusError();

  if(e.code != 1000) {
    err = "Server disconnected: " + e.code + ' ' + e.reason;
  }

  if(document.getElementById("loginSection") != null) {
	LoginActions.refresh(nickname, err);
  } else {
	LoginStore.setNickname(nickname);
	LoginStore.setError(err);
	React.render(
	 <LoginSection />,
	  document.getElementById("chattypantzapp")
	);
  }
};

// Sends a request to the server
ConnectionStore.chat.sendRequest = function(room, type, content) {
  this.socket.send(JSON.stringify({
    roomName: room,
    reqtype: type,
    content: content
  }));
};

// Register the store with the dispatcher.
ConnectionStore.dispatchToken = AppDispatcher.register(function(action) {
  var cs = ConnectionStore;
  var csc = cs.chat;
  switch (action.actionType) {
	case ActionTypes.REFRESH_CONNECTION:
		cs.emitChange();
		break;
	case ActionTypes.SEND_MESSAGE:
    	csc.sendRequest(csc.data.room, RequestTypes.MSG, action.message);
		cs.emitChange();
     	break;
    case ActionTypes.LOGOUT:
		csc.socket.close();
       break;
    default:
      // do nothing.
  }
});
