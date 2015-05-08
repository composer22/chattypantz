// TBD
var _messages = {};
var LoginStore = Object.assign({}, EventEmitter.prototype, {
  emitChange: function() {
    this.emit(EventTypes.CHANGE_EVENT);
  },

  addChangeListener: function(callback) {
    this.on(EventTypes.CHANGE_EVENT, callback);
  },

  removeChangeListener: function(callback) {
    this.removeListener(EventTypes.CHANGE_EVENT, callback);
  },

  get: function(id) {
    return _messages[id];
  },

  getAll: function(id) {
    return _messages;
  }

});

LoginStore.dispatchToken = AppDispatcher.register(function(action) {
  switch (action.actionType) {
    case ActionTypes.LOGIN:
		React.render(
			 <ConnectedSection />,
			  document.getElementById("chattypantzapp")
			);
      break;
    default:
      // do nothing.
  }
});
