// TBD
var LoginStore = Object.assign({}, EventEmitter.prototype, {

  nickname: "",
  error: "",

  emitChange: function() {
    this.emit(EventTypes.CHANGE_EVENT);
  },

  addChangeListener: function(callback) {
    this.on(EventTypes.CHANGE_EVENT, callback);
  },

  removeChangeListener: function(callback) {
    this.removeListener(EventTypes.CHANGE_EVENT, callback);
  },

  setNickname: function(nickname) {
	this.nickname = nickname;
  },

  getNickname: function() {
	return this.nickname;
  },

  setError: function(error) {
	this.error = error;
  },

  getError: function() {
	return this.error;
  }

});

LoginStore.dispatchToken = AppDispatcher.register(function(action) {
  var ls = LoginStore;
  var cs = ConnectionStore;
  switch (action.actionType) {
    case ActionTypes.LOGIN:
		ls.setNickname(action.nickname);
		cs.login();
		break;
    case ActionTypes.REFRESH_LOGIN:
		ls.setNickname(action.nickname);
		ls.setError(action.error);
		ls.emitChange();
      break;
    default:
      // do nothing.
  }
});
