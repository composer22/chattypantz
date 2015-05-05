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

	render: function() {
		return(
			<div>
			  <form className="ui error form" name="loginForm">
			    <div className="ui page grid">
			      <div className="three wide column">
			        <input type="submit" className="ui primary submit button" value="Start Chatting!" />
			      </div>
			      <div className="ten wide column field">
			        <input name="nickname" required type="text" placeholder="Enter your nickname and press the 'Start Chatting!' button..." />
			      </div>
			    </div>
			    <div className="ui page grid">
			      <div className="ui error message center aligned thirteen wide column hidden" id="connection-error">
			        <strong>Error:</strong>
			        <br/>
			        <br/> Please try again later...
			      </div>
			    </div>
			  </form>
			</div>
		);
	},

	componentDidUpdate: function() {


	}
});
