var React = require('react');

var selectedTabStore = require('../stores/selectedTab');
var privateChatActions = require('../actions/privateChat');
var tabActions = require('../actions/tab');

var UserListItem = React.createClass({
	handleClick() {
		var server = selectedTabStore.getServer();

		privateChatActions.open(server, this.props.user.nick);
		tabActions.select(server, this.props.user.nick);
	},

	render() {
		return <p onClick={this.handleClick}>{this.props.user.renderName}</p>;
	}
});

module.exports = UserListItem;