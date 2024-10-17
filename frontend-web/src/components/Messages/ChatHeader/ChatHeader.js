import React from 'react';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faEllipsisV } from '@fortawesome/free-solid-svg-icons';
import './ChatHeader.css';

const ChatHeader = ({ user, currentUserID }) => {

  const isGroup = user.conversation_type === 'Group';

  let avatarSrc;
  let nameDisplay;

  if (isGroup && user.users.length === 2) {
    const otherUser = user.users.find((u) => u.userID !== currentUserID);
    avatarSrc = otherUser ? otherUser.avatar : '../static/media/Logo.af2b2f1b32b135402e38.png';
    nameDisplay = otherUser ? otherUser.fullname : 'Unknown User';
  } else if (isGroup) {
    avatarSrc = user.avatar ? user.avatar : '../static/media/Logo.af2b2f1b32b135402e38.png';
    nameDisplay = user.name_conversation;
  } else {
    // Nếu là cuộc trò chuyện cá nhân
    const otherUser = user.users.find((u) => u.userID !== currentUserID);
    avatarSrc = otherUser ? otherUser.avatar : '../static/media/Logo.af2b2f1b32b135402e38.png';
    nameDisplay = otherUser ? otherUser.fullname : 'Unknown User';
  }

  return (
    <div className="chat-header " style={{ display: 'flex', justifyContent: 'space-between'}}>
      <div className="left-section">
        <div style={{ display: 'flex', alignItems: 'center' }}>
          <img
            src={avatarSrc}
            alt={nameDisplay}
            className="avatar"
            style={{width: '50px',height: '50px'}}
          />
          <div className="user-info">
            <h3>{nameDisplay}</h3>
            <span className="status">Online - Last seen, 2.02pm</span>
          </div>
        </div>

      </div>
      <div className="right-section">
        <div className="menu-icon ">
          <FontAwesomeIcon icon={faEllipsisV} />
        </div>
      </div>



    </div>
  );

};

export default ChatHeader;
