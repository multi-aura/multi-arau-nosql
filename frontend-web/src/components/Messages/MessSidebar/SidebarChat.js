import React, { useState } from 'react';

import MessageItem from '../MessageItem/MessageItem';
import './SidebarChat.css';
import { FaSearch } from 'react-icons/fa';

function SidebarChat({ conversations = [], onSelectChat }) {
  const [searchTerm, setSearchTerm] = useState('');
  const [filterType, setFilterType] = useState('All'); // All, Group, Single
  const [isSearchVisible, setSearchVisible] = useState(false); // State để kiểm soát hiển thị input

  // Lọc và tìm kiếm cuộc trò chuyện
  const filteredConversations = (conversations || [])
    .filter(conversation => {
      // Kiểm tra conversation và name_conversation tồn tại
      return conversation && conversation.name_conversation &&
        conversation.name_conversation.toLowerCase().includes(searchTerm.toLowerCase());
    })
    .filter(conversation => {
      if (filterType === 'Group') {
        return conversation.conversation_type === 'Group';
      } else if (filterType === 'Single') {
        return conversation.conversation_type === 'Single'; // Kiểm tra đúng với kiểu "Single"
      }
      return true;
    });

  return (
    <div className="sidebar-container">
      <div className="sidebar-header">
        <h5 style={{ color: "white" }}>Messages</h5>

        {!isSearchVisible && (
          <FaSearch
            className="search-icon"
            onClick={() => setSearchVisible(true)}
          />
        )}

        {isSearchVisible && (
          <input
            type="text"
            className="form-control search-input"
            placeholder="Search..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            // Ẩn khi người dùng nhấn Enter, nhưng giữ hiển thị trong suốt quá trình tìm kiếm
            onKeyDown={(e) => e.key === 'Enter' && setSearchVisible(false)}
          />
        )}

        <select
          className="form-select"
          value={filterType}
          onChange={(e) => setFilterType(e.target.value)}
        >
          <option value="All">All</option>
          <option value="Group">Group</option>
          <option value="Single">Single</option>
        </select>
      </div>

      <ul className="message-list">
        {filteredConversations.length > 0 ? (
          filteredConversations.map((conversation, index) => (
            <MessageItem
              key={index}
              message={conversation}
              onClick={() => onSelectChat(conversation._id)} // Đổi `message._id` thành `conversation._id`
            />
          ))
        ) : (
          <li className="no-messages">No messages found</li>
        )}
      </ul>
    </div>
  );
}

export default SidebarChat;
