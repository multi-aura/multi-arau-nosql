import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { getUserProfile } from '../services/usersService';
import UserInfo from '../components/UserProfile/UserInfo/UserInfo';
import Layout from '../layouts/Layout';
import Cookies from 'js-cookie';

const UserViewProfile = () => {
  const { username } = useParams();
  const [user, setUser] = useState(null);
  const [userData, setUserData] = useState(null);

  useEffect(() => {
    const fetchUserDetails = async () => {
      try {
        const userDataFromService = await getUserProfile(username);
        setUser(userDataFromService);
      } catch (error) {
        console.error('Lỗi khi lấy thông tin người dùng:', error);
      }
    };

    fetchUserDetails();

    const storedUser = localStorage.getItem('user');
    if (storedUser) {
      setUserData(JSON.parse(storedUser));
    } else {
      const authToken = Cookies.get('authToken');
      if (authToken) {
      }
    }
  }, [username]);

  if (!user || !userData) {
    return <p>Loading...</p>;
  }

  return (
    <Layout userData={userData}>
      <div className="user-profile-page container">
        <div className="row">
          <div className="col-md-8">
            <h1>{user.fullname}</h1> 
          </div>
          <div className="col-md-4">
            <UserInfo user={user} />
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default UserViewProfile;
