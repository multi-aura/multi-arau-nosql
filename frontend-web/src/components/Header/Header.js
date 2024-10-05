import React from 'react';
import './Header.css';

function Header() {
  return (
    <div className="header p-3 shadow-sm">
      <div className="d-flex justify-content-between align-items-center">
        <h4 className="mb-0 text-white">Feed</h4>
        <img src="" alt="User Avatar" className="rounded-circle avatar" />
      </div>
    </div>
  );
}

export default Header;
