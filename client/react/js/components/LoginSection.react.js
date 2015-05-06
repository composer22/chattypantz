/** @jsx React.DOM */

var LoginSection = React.createClass({

	getInitialState: function(){
		return {
			nickname: ""
		};
	},

	componentDidMount: function() {


	},

	componentWillUnmount: function(){

	},

	componentWillMount: function(){
	},

	render: function() {
		return(
			<div>
			  <form className="ui error form" name="loginForm" ng-submit="init()">
			    <div className="ui page grid">
			      <div className="three wide column">
			        <input type="submit" className="ui primary submit button disabled" ref="loginButton" value="Start Chatting!" onClick={this.handleSubmit}/>
			      </div>
			      <div className="ten wide column field">
			        <input name="nickname" required type="text" onChange={this.handleNicknameChange}
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

	componentDidUpdate: function() {


	},

	handleSubmit: function() {
		React.unmountComponentAtNode(document.getElementById("chattypantzapp"))
	},

	// When the input field changes, check the length and disable the button if it is empty.
	handleNicknameChange: function(event) {
		if(event.target.value.length > 0) {
			React.findDOMNode(this.refs.loginButton).className = "ui primary submit button";
		} else {
			React.findDOMNode(this.refs.loginButton).className = "ui primary submit button disabled";
		}
	}
});
