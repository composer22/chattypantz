// Helper functions for sending commands to the dispatcher from teh Login form.

var LoginActions = {
	login: function() {
		AppDispatcher.dispatch({
			actionType: ActionTypes.LOGIN
		});
	}
}
