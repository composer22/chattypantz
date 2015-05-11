// Helper functions for sending commands to the dispatcher from the Connection Form.

var ConnectionActions = {
	// Send a message to the server.
	send: function(message) {
		AppDispatcher.dispatch({
			actionType: ActionTypes.SEND_MESSAGE,
			message: message
		});
	},

	// Refresh the form elements
	refresh: function() {
		AppDispatcher.dispatch({
			actionType: ActionTypes.REFRESH_CONNECTION
		});
	},

	// Logout of the server.
	logout: function() {
		AppDispatcher.dispatch({
			actionType: ActionTypes.LOGOUT
		});
	}
}
