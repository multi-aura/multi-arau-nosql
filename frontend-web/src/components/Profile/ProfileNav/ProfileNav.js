import React from 'react';
import './ProfileNav.css';

function ProfileNav({ activeTab, onTabChange }) {
  return (
    <div className="profile-nav">
      <ul className="nav nav-tabs justify-content-center ">
        <li className={`nav-item ${activeTab === 'posts' ? 'active' : ''}`} onClick={() => onTabChange('posts')}>
          <a className="nav-link">Posts</a>
        </li>
        <li className={`nav-item ${activeTab === 'introduce' ? 'active' : ''}`} onClick={() => onTabChange('introduce')}>
          <a className="nav-link">Introduce</a>
        </li>
        <li className={`nav-item ${activeTab === 'friends' ? 'active' : ''}`} onClick={() => onTabChange('friends')}>
          <a className="nav-link">Friends</a>
        </li>
        <li className={`nav-item ${activeTab === 'images' ? 'active' : ''}`} onClick={() => onTabChange('images')}>
          <a className="nav-link">Images</a>
        </li>
        <li className={`nav-item ${activeTab === 'more' ? 'active' : ''}`} onClick={() => onTabChange('more')}>
          <a className="nav-link">More...</a>
        </li>

      </ul>
    </div>
  );
}

export default ProfileNav;
