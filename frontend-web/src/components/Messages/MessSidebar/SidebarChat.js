import React, { useState } from 'react';
import MessageItem from '../MessageItem/MessageItem';
import './SidebarChat.css';

function SidebarChat({ messages, onSelectChat }) {
  const [searchTerm, setSearchTerm] = useState('');

  const [filterType, setFilterType] = useState('All'); 

  const filteredMessages = messages
    .filter(message => message.name.toLowerCase().includes(searchTerm.toLowerCase()))
    .filter(message => {
      if (filterType === 'Group') {
        return message.type === 'group';
      } else if (filterType === 'Single') {
        return message.type === 'single'; 
      }
      return true; 
    });

  return (
    <div className="sidebar-container">
      <div className="sidebar-header">
        <h5>Messages</h5>
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
        {filteredMessages.map((message, index) => (
          <MessageItem 
            key={index} 
            message={message} 
            onClick={() => onSelectChat(message)} 
          />
        ))}
      </ul>
    </div>
  );
}

export default SidebarChat;
