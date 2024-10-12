import axios from 'axios';
import { API_URL } from '../config/config';
import Cookies from 'js-cookie';


export const getFriends = async () => {
    try {
      const token = Cookies.get('authToken');
      
      const response = await axios.get(`${API_URL}/relationships/friends`, {
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
  
      const response = await axios.get(`${API_URL}/relationships/followers`, {
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
  
      const response = await axios.get(`${API_URL}/relationships/followings`, {
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