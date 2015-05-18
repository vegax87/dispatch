var React = require('react');
var Reflux = require('reflux');
var _ = require('lodash');

var util = require('../util');
var searchStore = require('../stores/search');
var selectedTabStore = require('../stores/selectedTab');
var searchActions = require('../actions/search');
var PureMixin = require('../mixins/pure');

var Search = React.createClass({
	mixins: [
		PureMixin,
		Reflux.connect(searchStore, 'search'),
		Reflux.connect(selectedTabStore, 'selectedTab')
	],

	getInitialState() {
		return {
			search: searchStore.getState(),
			selectedTab: selectedTabStore.getState()
		};
	},

	componentDidUpdate(prevProps, prevState) {
		if (!prevState.search.get('show') && this.state.search.get('show')) {
			this.refs.input.getDOMNode().focus();
		}
	},

	handleChange(e) {
		var tab = this.state.selectedTab;

		if (tab.channel) {
			searchActions.search(tab.server, tab.channel, e.target.value);
		}
	},

	render() {
		var style = {
			display: this.state.search.get('show') ? 'block' : 'none'
		};

		var results = this.state.search.get('results').map(result => {
			return (
				<p key={result.id}>{util.timestamp(new Date(result.time * 1000))} {result.from} {result.content}</p>
			);
		});

		return (
			<div className="search" style={style}>
				<input 
					ref="input"
					className="search-input" 
					type="text"
					onChange={this.handleChange} />
				<div className="search-results">{results}</div>
			</div>
		);
	}
});

module.exports = Search;