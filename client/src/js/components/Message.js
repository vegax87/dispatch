import React, { PureComponent } from 'react';
import { timestamp, linkify } from '../util';

export default class Message extends PureComponent {
  handleSenderClick = () => {
    const { message, openPrivateChat, select } = this.props;

    openPrivateChat(message.server, message.from);
    select(message.server, message.from, true);
  };

  render() {
    const { message } = this.props;
    const classes = ['message'];
    let sender = null;

    if (message.type) {
      classes.push(`message-${message.type}`);
    }

    if (message.from) {
      sender = (
        <span>
          {' '}
          <span className="message-sender" onClick={this.handleSenderClick}>
            {message.from}
          </span>
        </span>
      );
    }

    const style = {
      paddingLeft: `${window.messageIndent + 15}px`,
      textIndent: `-${window.messageIndent}px`,
      ...this.props.style
    };

    return (
      <p className={classes.join(' ')} style={style}>
        <span className="message-time">{timestamp(message.time)}</span>
        {sender}
        <span>{' '}{linkify(message.message)}</span>
      </p>
    );
  }
}
