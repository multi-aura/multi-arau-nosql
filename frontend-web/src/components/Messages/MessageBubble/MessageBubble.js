import React from 'react';
import './MessageBubble.css';

const MessageBubble = ({ message, userAvatar }) => {

  // Chuyển đổi định dạng thời gian gửi tin nhắn
  const formattedTime = new Date(message.sender.added_at).toLocaleString('en-US', {
    hour: 'numeric',
    minute: 'numeric',
    hour12: true,
  });

  return (
    <div className={`messchat-item ${message.isSentByUser ? 'sent' : 'received'}`}>
      {!message.isSentByUser && (
        <img src={userAvatar} alt="sender-avatar" className="avatar" />
      )}
      <div className="message-content">
        <p className="message-text">{message.content.text}</p>
        <span className="message-time">{formattedTime}</span>
      </div>
    </div>
  );
};

export default MessageBubble;
