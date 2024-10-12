import React, { useEffect, useState } from 'react';
import '../assets/css/MyProfile.css';
import ProfileHeader from '../components/Profile/ProfileHeader/ProfileHeader';
import ProfileNav from '../components/Profile/ProfileNav/ProfileNav';
import Posts from '../components/Profile/SubProfile/MyPost'; 
import Introduce from '../components/Profile/SubProfile/Introduce';
import Friends from '../components/Profile/SubProfile/Friends';
import Layout from '../layouts/Layout';
import { useLocation, useNavigate } from 'react-router-dom';
import Cookies from 'js-cookie';

function MyProfile() {
  const [userData, setUserData] = useState(null);
  const [activeTab, setActiveTab] = useState('posts');
  const navigate = useNavigate();
  const location = useLocation();
  const authToken = Cookies.get('authToken');

  useEffect(() => {
    if (!authToken) {
      navigate('/');
    } else if (location.state && location.state.userData) {
      setUserData(location.state.userData);
      localStorage.setItem('user', JSON.stringify(location.state.userData));
    } else {
      const storedUser = localStorage.getItem('user');
      if (storedUser) {
        setUserData(JSON.parse(storedUser));
      }
    }
    const storedTab = localStorage.getItem('activeTab');
    if (storedTab) {
      setActiveTab(storedTab);
    }
  }, [authToken, location, navigate]);

  const handleTabChange = (tab) => {
    setActiveTab(tab);
    localStorage.setItem('activeTab', tab);
  };

  const renderContent = () => {
    switch (activeTab) {
      case 'posts':
        return <Posts />;
      case 'introduce':
        return <Introduce />;
      case 'friends':
        return <Friends />;
      case 'images':
        return 1;
      case 'more':
        return 2;
      default:
        return <Posts />;
    }
  };

  return (
    <Layout userData={userData}>
      <div className="container myprofile-page text-white py-5">
        <ProfileHeader userData={userData} />
        <ProfileNav activeTab={activeTab} onTabChange={handleTabChange} />
        <div className="row mt-4">
          <div className="col-md-12">
          {renderContent()}
          </div>
        </div>
      </div>
    </Layout>
  );
}

export default MyProfile;
