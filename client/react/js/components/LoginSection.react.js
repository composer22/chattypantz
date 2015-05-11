// LoginSection prompts the user for a nickname and navigates the initial connection to the server.
var LoginSection = React.createClass({
	render: function(){
		return(
			<div id="loginSection">
			  <form className="ui error form" name="loginForm">
			    <div className="ui page grid">
			      <div className="three wide column">
			        <input type="submit" className="ui primary submit button disabled" ref="loginButton" value="Start Chatting!" onClick={this._handleSubmit}/>
			      </div>
			      <div className="ten wide column field">
			        <input name="nickname" ref="nickname" required type="text" onChange={this._handleNicknameChange}
					    placeholder="Enter your nickname and press the 'Start Chatting!' button..." />
			      </div>
			    </div>
			    <div className="ui page grid">
			      <div className={this._errorClassCurrent()} ref="errorBox" id="connection-error">
			        <strong>Error:</strong>
			        <br/> {this.state.error}
			        <br/> Please try again later...
			      </div>
			    </div>
			  </form>
			</div>
		);
	},

	getInitialState: function() {
		return {
			nickname: LoginStore.getNickname(),
			error: LoginStore.getError()
		};
	},

	componentWillMount: function(){
		this.setState({
			nickname: LoginStore.getNickname(),
			error: LoginStore.getError()
		});
	},

	componentDidMount: function(){
		LoginStore.addChangeListener(this._onChange);
	},

	componentWillUnmount: function(){
		LoginStore.removeChangeListener(this._onChange);
	},

	// callback method for store to communicate any data change.
	_onChange: function() {
		this.setState({
			nickname: LoginStore.getNickname(),
			error: LoginStore.getError()
		});
	},

	_errorClassCurrent: function() {
		if(this.state.error == '') {
		 return "ui error message center aligned thirteen wide column hidden";
		}
		return "ui error message center aligned thirteen wide column";
	},

    // Submit button for login.
	_handleSubmit: function() {
		LoginActions.login(React.findDOMNode(this.refs.nickname).value);
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
