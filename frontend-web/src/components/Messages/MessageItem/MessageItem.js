import React from 'react';
import './MessageItem.css';

const MessageItem = ({ message, currentUserID, onClick }) => {
  // Kiểm tra xem đây là cuộc trò chuyện nhóm hay cá nhân
  const isGroup = message.conversation_type === 'Group';

  // Xử lý avatar và tên hiển thị
  let avatar;
  let nameDisplay;

  if (isGroup) {
    // Nếu là nhóm, hiển thị tên nhóm và avatar của nhóm
    avatar = message.avatar || '../static/media/Logo.af2b2f1b32b135402e38.png';
    nameDisplay = message.name_conversation || 'Unknown Group';
  } else {
    // Nếu là cuộc trò chuyện cá nhân, tìm người dùng còn lại (khác với currentUserID)
    const otherUser = message.users.find((user) => user.userID !== currentUserID);
    avatar = otherUser ? otherUser.avatar : '../static/media/Logo.af2b2f1b32b135402e38.png';
    nameDisplay = otherUser ? otherUser.fullname : 'Unknown User';
  }

  // Kiểm tra nội dung tin nhắn cuối cùng
  const lastMessage = message.chats && message.chats.length > 0
    ? (message.chats[message.chats.length - 1].content.text 
      ? message.chats[message.chats.length - 1].content.text 
      : message.chats[message.chats.length - 1].content.image 
      ? 'Hình ảnh' 
      : message.chats[message.chats.length - 1].content.voice_url
      ? 'Ghi âm' 
      : 'Bạn bè mới. Hãy gửi lời chào')
    : 'Bạn bè mới. Hãy gửi lời chào ';

  // Lấy thời gian của tin nhắn cuối cùng
  const lastMessageDate = message.chats && message.chats.length > 0
    ? new Date(message.chats[message.chats.length - 1].createdat)
    : null;

  // Hiển thị thời gian nếu tồn tại
  let timeDisplay = '';
  if (lastMessageDate) {
    const timeDiff = (new Date() - lastMessageDate) / (1000 * 60 * 60); // Tính số giờ
    if (timeDiff < 24) {
      timeDisplay = `${Math.floor(timeDiff)} hours ago`;
    } else {
      const days = Math.floor(timeDiff / 24);
      timeDisplay = `${days} day${days > 1 ? 's' : ''} ago`;
    }
  }

  return (
    <li className="list-group-item d-flex align-items-center message-item" style={{ borderBottom: "1px solid white" }} onClick={onClick}>
      <img src={avatar} alt="profile" className="avatar rounded-circle me-3" />

      <div className="message-info">
        <strong className="message-name">{nameDisplay}</strong>

        <div className="last-message-container" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <p className="small message-last" style={{ color: "#e9edef", marginRight: '10px' }}>
            {lastMessage}
          </p>

          {timeDisplay && (
            <span className="small" style={{ color: "#6f787d" }}>
              {timeDisplay}
            </span>
          )}
        </div>
      </div>
    </li>
  );
};

export default MessageItem;
