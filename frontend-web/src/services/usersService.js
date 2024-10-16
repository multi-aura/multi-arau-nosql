import axios from 'axios';
import { API_URL } from '../config/config';
import Cookies from 'js-cookie';
const USER_URL = `${API_URL}`;

export const getUserProfile = async (username) => {
  try {
    const token = Cookies.get('authToken');
    
    const response = await axios.get(`${USER_URL}/${username}`, {
      headers: {
        Authorization: `Bearer ${token}`, 
      },
    });
    
    return response.data.data;
  } catch (error) {
    console.error('Fail to load user profile:', error);
    throw error;
  }
};