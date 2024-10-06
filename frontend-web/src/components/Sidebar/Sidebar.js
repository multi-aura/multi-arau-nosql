import React from 'react';
import './Sidebar.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faHome, faSearch, faCommentDots, faBell, faUser } from '@fortawesome/free-solid-svg-icons'; 

function Sidebar() {
  return (
    <div className="sidebar col-2 ">
      <h2 className="text-center py-4">Multi Aura</h2>
      <ul className="nav flex-column">
        <li className="nav-item">
          <a className="tab-link" href="#home">
            <FontAwesomeIcon icon={faHome} /> Home
          </a>
        </li>
        <li className="nav-item">
          <a className="tab-link" href="#search">
            <FontAwesomeIcon icon={faSearch} /> Search
          </a>
        </li>
        <li className="nav-item">
          <a className="tab-link" href="#messages">
            <FontAwesomeIcon icon={faCommentDots} /> Messages
          </a>
        </li>
        <li className="nav-item">
          <a className="tab-link" href="#notifications">
            <FontAwesomeIcon icon={faBell} /> Notifications
          </a>
        </li>
        <li className="nav-item">
          <a className="tab-link" href="#profile">
            <FontAwesomeIcon icon={faUser} /> Profile
          </a>
        </li>
      </ul>
  </div>
  );
}

export default Sidebar;