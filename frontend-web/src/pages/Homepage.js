import React,{ useEffect, useState } from 'react';
import Layout from '../layouts/Layout';
import Feed from '../components/Feed/Feed';
import useSession from '../hooks/useSession';
import { useLocation, useNavigate } from 'react-router-dom';
import Cookies from 'js-cookie'; 

function Homepage() {
    const [userData, setUserData] = useState(null);
   
    const navigate = useNavigate();
    const location = useLocation();

    const authToken = Cookies.get('authToken'); 
    useEffect(() => {
      if (!authToken) {
        navigate('/');
      } else if (location.state && location.state.userData) {
        setUserData(location.state.userData);
      }
    }, [authToken, location, navigate]);
    
    return (
      <Layout userData={userData}>
       <Feed />
      </Layout>
    );
}
export default Homepage;