import React from 'react';
import './ProfileHeader.css'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCog } from '@fortawesome/free-solid-svg-icons';
function ProfileHeader({ userData }) {
    if (!userData || !userData.avatar) {
        console.log(userData);
        return <p>Loading profile...</p>;  
      }
  return (
    <div className="row align-items-center my-4">
      <div className="col-md-4 text-center">
        <img src={userData.avatar} alt="Avatar" className="rounded-circle profile-avatar" />
      </div>
      <div className="col-md-8 profile-info">
        <h2>{userData.fullname}</h2>
        <p>50 posts • 1.1k friends • 20k likes</p>
        {/* <p>{userData.posts.length} posts • {userData.friends} friends • {userData.likes} likes</p> */}
        <button className="btn btn-outline-light d-flex align-items-center">
          <FontAwesomeIcon icon={faCog} className="me-2" />
          Edit Profile
        </button>
      </div>
    </div>
  );
}

export default ProfileHeader;
