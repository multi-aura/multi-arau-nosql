import React from 'react';
import './MessageItem.css';

const MessageItem = ({ message, onClick }) => {
  return (
    <li className="list-group-item text d-flex align-items-center message-item" onClick={onClick}>
      <img src={message.avatar} alt="profile" className="avatar me-2" />
      <div className="message-info">
        <strong className="message-name">{message.name}</strong>
        <p className="small message-last">{message.lastMessage}</p>
      </div>
    </li>
  );
};

export default MessageItem;
