// TBD
var ConnectionStore = Object.assign({}, EventEmitter.prototype, {
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
    return null;
  },

  getAll: function(id) {
    return null;
  }

});

ConnectionStore.dispatchToken = AppDispatcher.register(function(action) {
  switch (action.actionType) {
	case ActionTypes.SEND_MESSAGE:
		console.log(action.message);
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
