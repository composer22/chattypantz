// LoginSection prompts the user for a nickname and navigates the initial connection to the server.
var LoginSection = React.createClass({
	// setState(function|object nextState[, function callback])
	// forceUpdate([function callback])

	// object propTypes
	// array mixins
	// object statics
	// string displayName

	// Built-in Component Methods

	render: function(){
		return(
			<div>
			  <form className="ui error form" name="loginForm">
			    <div className="ui page grid">
			      <div className="three wide column">
			        <input type="submit" className="ui primary submit button disabled" ref="loginButton" value="Start Chatting!" onClick={this._handleSubmit}/>
			      </div>
			      <div className="ten wide column field">
			        <input name="nickname" required type="text" onChange={this._handleNicknameChange}
					    placeholder="Enter your nickname and press the 'Start Chatting!' button..." />
			      </div>
			    </div>
			    <div className="ui page grid">
			      <div className="ui error message center aligned thirteen wide column hidden"
				         ng-show="chat.status.error" id="connection-error">
			        <strong>Error:</strong>
			        <br/> chat.status.error
			        <br/> Please try again later...
			      </div>
			    </div>
			  </form>
			</div>
		);
	},


	getInitialState: function() {
		return {
			nickname: ""
		};
	},

	getDefaultProps: function() {
		// NOP
	},

	// Built-in Lifecycle Methods

	componentWillMount: function(){
		// NOP
	},

	componentDidMount: function(){
		LoginStore.addChangeListener(this._onChange);
	},

	componentWillReceiveProps: function(){
		// NOP
	},
	shouldComponentUpdate: function(){
		// NOP
	},
	componentWillUpdate: function(){
		// NOP
	},
	componentDidUpdate: function(){
		// NOP
	},

	componentWillUnmount: function(){
		LoginStore.removeChangeListener(this._onChange);
	},

	// Custom Methods

	// callback method for store to communicate any data change.
	_onChange: function() {
		// on call, you would pull data from the LoginStore for display
	},

    // Submit button for login.
	_handleSubmit: function() {
		LoginActions.login();
		return false;
	},

	// When the input field changes, check the length and disable the button if it is empty.
	_handleNicknameChange: function(event) {
		if(event.target.value.length > 0) {
			React.findDOMNode(this.refs.loginButton).className = "ui primary submit button";
		} else {
			React.findDOMNode(this.refs.loginButton).className = "ui primary submit button disabled";
		}
	}
});