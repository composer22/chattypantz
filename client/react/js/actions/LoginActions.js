// Helper functions for sending commands to the dispatcher from teh Login form.

var LoginActions = {
	login: function(nickname) {
		AppDispatcher.dispatch({
			actionType: ActionTypes.LOGIN,
			nickname: nickname
		});
	},
	refresh: function(nickname, error) {
		AppDispatcher.dispatch({
			actionType: ActionTypes.REFRESH_LOGIN,
			nickname: nickname,
			error: error
		});
	}
}
