import React from 'react';
import './SuggestedUser.css';

const SuggestedUser = ({ user }) => {
  return (
    <li key={user.id} className="list-group-item suggested-user d-flex justify-content-between align-items-center">
      <div className="d-flex align-items-center">
        <img src={user.avatar} alt={user.name} className="rounded-circle me-3 user-avatar" />
        <div className="user-info">
          <p className="mb-0 user-name">{user.name}</p>
        </div>
      </div>
      <button className="btn btn-outline-light btn-sm follow-btn">Follow</button>
    </li>
  );
};

export default SuggestedUser;
