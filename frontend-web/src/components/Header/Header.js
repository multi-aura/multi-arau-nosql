import React from 'react';
import './Header.css';

function Header({ userData }) {
  return (
    <div className="header p-3 shadow-sm">
      <div className="d-flex justify-content-between align-items-center">
        <h4 className="mb-0 text-white">Feed</h4>
        {userData ? (  
          <div className="d-flex align-items-center">
            <span className="text-white me-3">{userData.fullname}</span>
            <img src={userData.avatar} alt="User Avatar" className= "rounded-circle avatar" />
          </div>
        ) : (
          <span className="text-white">Guest</span>  
        )}
      </div>
    </div>
  );
}
export default Header;
