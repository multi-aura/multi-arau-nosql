import React, { useState } from 'react';
import { followUser } from '../../../services/RelationshipService'; 
import './SuggestedUser.css';

const SuggestedUser = ({ user }) => {
  const [isFollowing, setIsFollowing] = useState(user.isFollowing);

  const handleFollowClick = async () => {
    try {
      await followUser(user.userID); 
      setIsFollowing(true);  
    } catch (error) {
      console.error('Failed to follow user:', error);
    }
  };
  return (
    <li key={user.id} className="list-group-item suggested-user d-flex justify-content-between align-items-center">
      <div className="d-flex align-items-center">
        <img src={user.avatar} alt={user.username} className="rounded-circle me-3 user-avatar" />
        <div className="user-info">
          <p className="mb-0 user-name">{user.fullname}</p>
        </div>
      </div>
      <button 
        className="follow-btn" 
        onClick={handleFollowClick} 
        disabled={isFollowing}
      >
        {isFollowing ? 'Following' : 'Follow'}
      </button>
    </li>
  );
};

export default SuggestedUser;
