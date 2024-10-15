import React from 'react';
import './PersonItem.css'; 

const PersonItem = ({ person }) => {
  return (
    <li className="person-item">
      <img src={person.avatar} alt={person.username} className="people-avatar" />
      <div className="person-info">
        <p className="person-name">{person.fullname}</p>
        <p className="person-username">{person.username}</p>
      </div>
      <button className="follow-btn">{person.isFollowing ? 'Following' : 'Follow'}</button>
    </li>
  );
};

export default PersonItem;
