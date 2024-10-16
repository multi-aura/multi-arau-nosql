import React from 'react';
import './PersonItem.css'; 

const PersonItem = ({ person }) => {
  return (
    <li className="person-item">
      <img src={person.avatar || 'https://firebasestorage.googleapis.com/v0/b/multi-aura.appspot.com/o/Hihon%2F1728534046_9ea1c9841cadbef3e7bc.jpg?alt=media&token=3d221a08-d064-4ece-881a-32e2c5d273e1'} alt={person.username} className="people-avatar" />
      <div className="persons-info">
        <p className="persons-name">{person.fullname}</p>
        <p className="persons-username">{person.username}</p>
      </div>
      <button className="follow-btn">{person.isFollowing ? 'Following' : 'Follow'}</button>
    </li>
  );
};

export default PersonItem;
