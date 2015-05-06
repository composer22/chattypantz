/** @jsx React.DOM */

var ConnectedSection = React.createClass({

	getInitialState: function(){
		return {
			nickname: ""
		};
	},

	componentDidMount: function() {

	},

	componentWillUnmount: function(){

	},

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

		  <form className="ui form" name="sendForm" ng-submit="sendButtonClicked()">
		    <div className="ui page grid">
		      <div className="two wide column">
		        <input type="submit" className="ui submit button primary" data-ng-disabled="chat.data.messageText == ''" value="Send!" />
		      </div>
		      <div className="eight wide column field">
		        <input className="message" ng-model="chat.data.messageText" name="messageText" type="text" placeholder="Type message here..." />
		      </div>
		      <div className="two wide column">
		        <input type="button" class="ui mini red button" ng-click="quit()" value="Quit" />
		      </div>
		    </div>
		  </form>
		</div>
		);
	},

	componentDidUpdate: function() {


	}
});
