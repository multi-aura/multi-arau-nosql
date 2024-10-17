import axios from 'axios';
import { API_URL } from '../config/config';
import Cookies from 'js-cookie';

const RELATIONSHIPS_URL = `${API_URL}/relationships`;
const FOLLOW_URL = `${RELATIONSHIPS_URL}/follow`;
export const getFriends = async () => {
    try {
      const token = Cookies.get('authToken');
      
      const response = await axios.get(`${RELATIONSHIPS_URL}/friends`, {
        headers: {
          Authorization: `Bearer ${token}`, 
        },
      });
  
      return response.data.data; 
    } catch (error) {
      console.error('Fail to load list friend:', error);
      throw error; 
    }
};
export const getFollowers = async () => {
    try {
      const token = Cookies.get('authToken');
  
      const response = await axios.get(`${RELATIONSHIPS_URL}/followers`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
  
      return response.data.data; 
    } catch (error) {
      console.error('Fail to load list follower:', error);
      throw error;
    }
};
export const getFollowings = async () => {
    try {
      const token = Cookies.get('authToken');
  
      const response = await axios.get(`${RELATIONSHIPS_URL}/followings`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
  
      return response.data.data;
    } catch (error) {
      console.error('Fail to load following:', error);
      throw error;
    }
};

export const unfollowUser = async (userID) => {
  try {
      const token = Cookies.get('authToken');
      
      const response = await axios.delete(`${RELATIONSHIPS_URL}/unfollow/${userID}`, {
          headers: {
              Authorization: `Bearer ${token}`,
          },
      });

      return response.data; 
  } catch (error) {
      console.error(`Fail to unfollow user ${userID}:`, error);
      throw error;
  }
};

export const followUser = async (userID) => {
  try {
      const token = Cookies.get('authToken');
      if (!token) {
          throw new Error('Token không tồn tại. Vui lòng đăng nhập.');
      }

      const response = await axios.post(`${FOLLOW_URL}/${userID}`, {}, {
          headers: {
              Authorization: `Bearer ${token}`,
          },
      });
      
      return response.data;
  } catch (error) {
      console.error('Fail to follow this user:', error.response ? error.response.data : error.message);
      throw error;
  }
};
export const checkRelationshipStatus = async (userID) => {
  try {
    const token = Cookies.get('authToken');

    const response = await axios.post(`${RELATIONSHIPS_URL}/status/${userID}`, {}, {
      headers: {
        Authorization: `Bearer ${token}`, 
      },
    });
    
    return response.data.data; 
  } catch (error) {
    console.error('Fail to check Relationship:', error);
    throw error;
  }
};
export const blockUser = async (userID) => {
  const token = Cookies.get('authToken');
  const response = await axios.post(`${RELATIONSHIPS_URL}/block/${userID}`, {}, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return response.data;
};

export const unblockUser = async (userID) => {
  const token = Cookies.get('authToken');
  const response = await axios.post(`${RELATIONSHIPS_URL}/unblock/${userID}`, {}, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return response.data;
};