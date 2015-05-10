// Tracks information about the connection to the server and the display

var ConnectionStore = Object.assign({}, EventEmitter.prototype, {
   chat: {
    socket: null,
    status: {
      error: null,
      starting: false,
      started: false
    },
    data: {
      users: [],
      room: 'Demo',
      nickname: '',
      history: '',
      messageText: ''
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

  getSocketError: function() {
	 return this.chat.status.error;
  }

});

ConnectionStore.chat.init = function() {
    // Open a socket and register event handlers.
    this.socket = new WebSocket(Server);
    this.socket.onopen = this.onOpenWS;
    this.socket.onmessage = this.onMessageWS;
    this.socket.onerror = this.onErrorWS;
    this.socket.onclose = this.onErrorWS;
    this.status.starting = true;
};

//////// SOCKET EVENT HANDLERS ////////

ConnectionStore.chat.onOpenWS = function(e) {
	var csc = ConnectionStore.chat;
    var ls = LoginStore;
    csc.status.starting = false;
    csc.status.started = true;
    csc.status.error = null;
    csc.data.nickname = ls.getNickname();
    // Set the nickname, join the demo room, and get a list of names.
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
      csc.data.history += "Chattypantz server: " + response.content + '\n';
      csc.data.history += "Quitting Chattypantz.\n";
      break;
    default:
      csc.data.history += response.content + '\n';
      csc.data.history += "Quitting Chattypantz.\n";
  }
};

// Connection error.
ConnectionStore.chat.onErrorWS = function(e) {
  var csc = ConnectionStore.chat;
  var ls = LoginStore;
  // Stop run and show err.
  csc.status.error = "Server disconnected: " + e.reason;
  csc.status.started = false;
  ls.setError(csc.status.error);
  if(typeof LoginSection != "undefined") {
	LoginActions.refresh(csc.data.nickname, csc.status.error);
  } else {
	ls.setNickname('');
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
  switch (action.actionType) {
	case ActionTypes.SEND_MESSAGE:
     	break;
    case ActionTypes.LOGOUT:
		React.render(
			 <LoginSection />,
			  document.getElementById("chattypantzapp")
			);
       break;
    default:
      // do nothing.
  }
});
