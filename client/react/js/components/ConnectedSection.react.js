// ConnectedSection is the interactive form for communication to the server once connected.
var ConnectedSection = React.createClass({
	render: function() {
		return(
		<div id="connectedSection">
		  <div className="ui page grid">
		    <div className="twelve wide column">
		      <textarea id="chat-history" ref="chatHistory"></textarea>
		    </div>
		    <div className="four wide column" id="online-user-list">
		      <a className="ui teal ribbon label">Online Users {this.state.users.length}</a>
		      <br/>
		       {this._displayUsers()}
		    </div>
		  </div>

		  <form className="ui form" name="sendForm">
		    <div className="ui page grid">
		      <div className="two wide column">
		        <input type="submit" className="ui primary submit button disabled" ref="sendButton" value="Send!" onClick={this._handleSend}/>
		      </div>
		      <div className="eight wide column field">
		        <input className="message" name="messageText" type="text" onChange={this._handleMessageChange}
				   ref="messageBox" placeholder="Type message here..." />
		      </div>
		      <div className="two wide column">
		        <input type="button" className="ui mini red button" value="Quit" onClick={this._handleQuit}/>
		      </div>
		    </div>
		  </form>
		</div>
		);
	},

	getInitialState: function() {
		return {
			users: [],
           history: ''
		};
	},

	componentDidMount: function(){
		ConnectionStore.addChangeListener(this._onChange);
	},

	componentWillUnmount: function(){
		ConnectionStore.removeChangeListener(this._onChange);
	},

	_disableSendButton: function() {
		React.findDOMNode(this.refs.sendButton).className = "ui primary submit button disabled";
	},

	_enableSendButton: function() {
		React.findDOMNode(this.refs.sendButton).className = "ui primary submit button";
	},

	// callback method for store to communicate any data change.
	_onChange: function() {
		var csc = ConnectionStore.chat;
		this.setState({
			users: csc.data.users,
			history: csc.data.history
		});
		var ch = React.findDOMNode(this.refs.chatHistory);
        ch.value = this.state.history;
        ch.scrollTop = ch.scrollHeight;
		this.forceUpdate();
	},

	_displayUsers: function() {
		var result = ""
		return (
		<div>
			{this.state.users.map(function(user) {
				return <span>{user}<br/></span>;
			})}
		</div>
		);
	},

	// Send button pressed for message.
	_handleSend: function() {
		var msg = React.findDOMNode(this.refs.messageBox).value
		ConnectionActions.send(msg);
		React.findDOMNode(this.refs.messageBox).value = "";
		this._disableSendButton();
		return false;
	},

	// When the message field changes, check the length and disable the
	// send button if it is empty.
	_handleMessageChange: function(event) {
		if(event.target.value.length > 0) {
			this._enableSendButton();
		} else {
			this._disableSendButton();
		}
	},

	// Quit button pressed, so disconnect; return to main page.
	_handleQuit: function() {
		ConnectionActions.logout();
		return false;
	}
});
