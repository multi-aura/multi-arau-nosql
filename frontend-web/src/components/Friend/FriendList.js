import React from 'react';
import FriendItem from './FriendItem';
import './FriendList.css'
const FriendList = ({ friends }) => {
  if (friends.length === 0) {
    return <p>Không có bạn bè</p>;
  }

  return (
    <ul className="list-group friend-list-container">
      {friends.map((friend) => (
        <FriendItem key={friend.userID} friend={friend} />
      ))}
    </ul>
  );
};

export default FriendList;
