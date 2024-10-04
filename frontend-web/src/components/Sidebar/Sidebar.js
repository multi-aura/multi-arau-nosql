import React from 'react';
import './Sidebar.css';

function Sidebar() {
  return (
    <div className="sidebar">
      <h2>Multi Aura</h2>
      <ul>
        <li>Home</li>
        <li>Search</li>
        <li>Messages</li>
        <li>Notifications</li>
        <li>Profile</li>
      </ul>
    </div>
  );
}

export default Sidebar;
