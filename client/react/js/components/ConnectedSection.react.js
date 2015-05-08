// ConnectedSection is the interactive form for communication to the server once connected.
var ConnectedSection = React.createClass({
	// setState(function|object nextState[, function callback])
	// forceUpdate([function callback])

	// object propTypes
	// array mixins
	// object statics
	// string displayName

	// Built-in Component Methods

	render: function() {
		return(
		<div>
		  <div className="ui page grid">
		    <div className="twelve wide column">
		      <textarea id="chat-history"></textarea>
		    </div>
		    <div className="four wide column" id="online-user-list">
		      <a className="ui teal ribbon label">Online Users (chat.data.users.length)</a>
		      <br/>
		      <span ng-repeat="nickname in chat.data.users">nickname
		        <br/>
		      </span>
		    </div>
		  </div>

		  <form className="ui form" name="sendForm">
		    <div className="ui page grid">
		      <div className="two wide column">
		        <input type="submit" className="ui primary submit button disabled" ref="sendButton" value="Send!" onClick={this._handleSend}/>
		      </div>
		      <div className="eight wide column field">
		        <input className="message" name="messageText" type="text" onChange={this._handleMessageChange}
				   placeholder="Type message here..." />
		      </div>
		      <div className="two wide column">
		        <input type="button" class="ui mini red button" value="Quit" onClick={this._handleQuit}/>
		      </div>
		    </div>
		  </form>
		</div>
		);
	},

	getInitialState: function() {
		return null;
	},
	getDefaultProps: function() {
		// NOP
	},

	// Built-in Lifecycle Methods

	componentWillMount: function(){
		// NOP
	},
	componentDidMount: function(){
		ConnectionStore.addChangeListener(this._onChange);
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
		ConnectionStore.removeChangeListener(this._onChange);
	},

	// callback method for store to communicate any data change.
	_onChange: function() {
		// on call, you would pull data from the ConnectionStore for display
	},

	// Send button pressed for message.
	_handleSend: function() {
		ConnectionActions.send("TODO The Message Goes Here.");
		return false;
	},

	// When the message field changes, check the length and disable the send button if it is empty.
	_handleMessageChange: function(event) {
		if(event.target.value.length > 0) {
			React.findDOMNode(this.refs.sendButton).className = "ui primary submit button";
		} else {
			React.findDOMNode(this.refs.sendButton).className = "ui primary submit button disabled";
		}
	},

	// Quit button pressed, so disconnect; return to main page.
	_handleQuit: function() {
		ConnectionActions.logout();
		return false;
	}
});
