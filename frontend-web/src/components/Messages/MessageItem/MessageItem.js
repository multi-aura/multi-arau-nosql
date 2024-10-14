import React from 'react';
import './MessageItem.css';

const MessageItem = ({ message, onClick }) => {
  const user = message.users && message.users.length > 0 ? message.users[0] : null;
  const avatar = user && user.avatar ? user.avatar : 'https://www.w3schools.com/w3images/avatar2.png';

  const lastMessage = message.chats && message.chats.length > 0
    ? message.chats[message.chats.length - 1].content.text ||
    (message.chats[message.chats.length - 1].content.image ? 'Sent an image' : 'Sent a voice message')
    : 'No recent messages';

  // Lấy thời gian của tin nhắn cuối cùng
  const lastMessageDate = message.chats && message.chats.length > 0
    ? new Date(message.chats[message.chats.length - 1].createdat)
    : null;

  // Tính khoảng cách thời gian hiện tại với tin nhắn cuối
  const now = new Date();
  const timeDiff = lastMessageDate ? (now - lastMessageDate) / (1000 * 60 * 60) : null; // Tính số giờ

  // Logic hiển thị thời gian
  let timeDisplay = '';
  if (timeDiff) {
    if (timeDiff < 24) {
      timeDisplay = `${Math.floor(timeDiff)} hours ago`; // Hiển thị số giờ
    } else {
      const days = Math.floor(timeDiff / 24);
      timeDisplay = `${days} day${days > 1 ? 's' : ''} ago`; // Hiển thị số ngày
    }
  }

  return (
    <li className="list-group-item d-flex align-items-center message-item" style={{ borderBottom: "1px solid white" }} onClick={onClick}>
      {/* Avatar của người dùng */}
      <img src={avatar} alt="profile" className="avatar rounded-circle me-3" />

      {/* Thông tin cuộc trò chuyện */}
      <div className="message-info">
        <strong className="message-name">{message.name_conversation || (user && user.fullname)}</strong>

        {/* Tin nhắn và thời gian trên cùng một dòng */}
        <div className="last-message-container" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <p className="small message-last" style={{ color: "#e9edef", marginRight: '10px' }}>
            {lastMessage}
          </p>

          {/* Hiển thị thời gian tin nhắn cuối cùng */}
          {timeDisplay && (
            <span className="small" style={{color:"#6f787d"}}>
              {timeDisplay}
            </span>
          )}
        </div>
      </div>


    </li>
  );
};

export default MessageItem;
