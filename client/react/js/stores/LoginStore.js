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
  switch (action.actionType) {
    case ActionTypes.LOGIN:
		LoginStore.setNickname(action.nickname);
		ConnectionStore.login();
		break;
    case ActionTypes.REFRESH_LOGIN:
		LoginStore.setNickname(action.nickname);
		LoginStore.setError(action.error);
		LoginStore.emitChange();
      break;
    default:
      // do nothing.
  }
});
