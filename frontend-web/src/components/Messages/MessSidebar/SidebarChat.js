import React, { useState, useEffect } from 'react';
import MessageItem from '../MessageItem/MessageItem';
import './SidebarChat.css';
import { FaSearch } from 'react-icons/fa';

function SidebarChat({ conversations = [], onSelectChat, newMessageItems }) {
  console.log(conversations);

  const [searchTerm, setSearchTerm] = useState('');
  const [filterType, setFilterType] = useState('All'); // All, Group, Single
  const [isSearchVisible, setSearchVisible] = useState(false); // State để kiểm soát hiển thị input
  const [filteredConversations, setFilteredConversations] = useState([]);
  const [sortedConversations, setSortedConversations] = useState([]);

  // Sắp xếp các cuộc trò chuyện dựa trên tin nhắn mới nhất (updatedat)
  useEffect(() => {
    if (conversations.length > 0) {
      const sorted = conversations.slice().sort((a, b) => {
        const lastMessageA = a.chats.length > 0 ? a.chats[a.chats.length - 1].updatedat : '0000-00-00T00:00:00Z';
        const lastMessageB = b.chats.length > 0 ? b.chats[b.chats.length - 1].updatedat : '0000-00-00T00:00:00Z';
        return new Date(lastMessageB) - new Date(lastMessageA);
      });
      setSortedConversations(sorted);
    }
  }, [conversations]);

  // Lọc và tìm kiếm cuộc trò chuyện dựa trên danh sách đã sắp xếp
  useEffect(() => {
    const filtered = sortedConversations
      .filter(conversation => {
        return conversation && conversation.name_conversation &&
          conversation.name_conversation.toLowerCase().includes(searchTerm.toLowerCase());
      })
      .filter(conversation => {
        if (filterType === 'Group') {
          return conversation.conversation_type === 'Group';
        } else if (filterType === 'Single') {
          return conversation.conversation_type === 'Single';
        }
        return true;
      });

    setFilteredConversations(filtered);
  }, [sortedConversations, searchTerm, filterType]);

  // Cập nhật tin nhắn mới vào cuộc trò chuyện và sắp xếp lại
  useEffect(() => {
    if (newMessageItems) {
      setSortedConversations((prevConversations) => {
        const updatedConversations = prevConversations.map(conversation => {
          if (conversation._id === newMessageItems.conversationID) {
            // Cập nhật tin nhắn và thời gian tin nhắn cuối cùng
            return {
              ...conversation,
              lastMessage: newMessageItems.content.text || "",
              lastMessageTime: newMessageItems.createdat || new Date().toISOString(),
              chats: [...conversation.chats, { ...newMessageItems, updatedat: newMessageItems.createdat }]
            };
          }
          return conversation;
        });

        // Sắp xếp lại danh sách dựa trên tin nhắn mới nhất
        const sorted = updatedConversations.slice().sort((a, b) => {
          const lastMessageA = a.chats.length > 0 ? a.chats[a.chats.length - 1].updatedat : '0000-00-00T00:00:00Z';
          const lastMessageB = b.chats.length > 0 ? b.chats[b.chats.length - 1].updatedat : '0000-00-00T00:00:00Z';
          return new Date(lastMessageB) - new Date(lastMessageA);
        });
        return sorted;
      });
    }
  }, [newMessageItems]);

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
              onClick={() => onSelectChat(conversation._id)}
            />
          ))
        ) : (
          <li className="no-messages">
            <span role="img" aria-label="no-messages" className="no-messages-icon">📭</span>
            No messages found
          </li>
        )}
      </ul>

    </div>
  );
}

export default SidebarChat;
