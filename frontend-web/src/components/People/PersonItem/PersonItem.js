import React from 'react';
import './PersonItem.css'; 

const PersonItem = ({ person }) => {
  return (
    <li className="person-item">
      <img src={person.avatar} alt={person.name} className="people-avatar" />
      <div className="person-info">
        <p className="person-name">{person.name}</p>
        <p className="person-description">{person.description}</p>
      </div>
      <button className="follow-btn">{person.isFollowing ? 'Following' : 'Follow'}</button>
    </li>
  );
};

export default PersonItem;
