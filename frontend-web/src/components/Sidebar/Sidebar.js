import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import './Sidebar.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faHome, faThLarge , faCommentDots, faBell, faUser } from '@fortawesome/free-solid-svg-icons'; 

function Sidebar() {
  const location = useLocation();
  const [activeTab, setActiveTab] = useState('/Home'); 
  useEffect(() => {
    setActiveTab(location.pathname);
  }, [location.pathname]);
  const handleTabClick = (tab) => {
    setActiveTab(tab);
  };

  return (
    <div className="sidebar">
      <h2 className="text-center py-4">Multi Aura</h2>
      <ul className="nav flex-column">
        <li className="nav-item">
          <a 
            className={`tab-link ${activeTab === '/Home' ? 'active' : ''}`} 
            href="/Home" 
            onClick={() => handleTabClick('/Home')}
          >
            <FontAwesomeIcon icon={faHome} /> Home
          </a>
        </li>
        <li className="nav-item">
          <a 
            className={`tab-link ${activeTab === '#search' ? 'active' : ''}`} 
            href="/explore" 
            onClick={() => handleTabClick('#search')}
          >
            <FontAwesomeIcon icon={faThLarge} /> Explore
          </a>
        </li>
        <li className="nav-item">
          <a 
            className={`tab-link ${activeTab === '#messages' ? 'active' : ''}`} 
            href="/chat" 
            onClick={() => handleTabClick('#messages')}
          >
            <FontAwesomeIcon icon={faCommentDots} /> Messages
          </a>
        </li>
        <li className="nav-item">
          <a 
            className={`tab-link ${activeTab === '#notifications' ? 'active' : ''}`} 
            href="#notifications" 
            onClick={() => handleTabClick('#notifications')}
          >
            <FontAwesomeIcon icon={faBell} /> Notifications
          </a>
        </li>
        <li className="nav-item">
          <a 
            className={`tab-link ${activeTab === '/profile' ? 'active' : ''}`} 
            href="/profile" 
            onClick={() => handleTabClick('/profile')}
          >
            <FontAwesomeIcon icon={faUser} /> Profile
          </a>
        </li>
      </ul>
  </div>
  );
}

export default Sidebar;
