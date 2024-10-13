import axios from 'axios';
import { API_URL } from '../config/config';
import Cookies from 'js-cookie';

const RELATIONSHIPS_URL = `${API_URL}/relationships`;
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
