import React, { useState }  from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'; 
import { faUserFriends, faUserMinus   } from '@fortawesome/free-solid-svg-icons';
import './FriendItem.css'
const FriendItem = ({ friend }) => {
  const [isFriend, setIsFriend] = useState(true);
 

  const handleUnfriend = () => {
    setIsFriend(false); 
  };


  return (
    <li className="list-group-item friend-item d-flex justify-content-between align-items-center py-3">
      <div className="d-flex align-items-center">
        <img
          src={friend.avatar || 'https://via.placeholder.com/50'}
          alt={friend.fullname}
          className="rounded-circle me-3"
          style={{ width: '60px', height: '60px', objectFit: 'cover' }}
        />
        <div>
          <h5 className="mb-0">{friend.fullname}</h5>
          <small className="text">{friend.username}</small>
        </div>  
      </div>
      <button className="btn btn-dark d-flex align-items-center" onClick={handleUnfriend}>
        {isFriend ? 'Friend' : 'Unfriend'}
        <FontAwesomeIcon 
          icon={isFriend ? faUserFriends : faUserMinus}
          className="fa-user-friends ms-2"
        />
      </button>
      
    </li>
  );
};

export default FriendItem;
