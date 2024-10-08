import axios from 'axios';
import { API_URL } from '../config/config';
import Cookies from 'js-cookie';
export const login = async (username, password) => {
  try {
    const response = await axios.post(`${API_URL}/user/login`, {
      username,
      password,
    });
    const { token } = response.data.data; 
    Cookies.set('authToken', token, {
      expires: 1, 
      secure: process.env.NODE_ENV === 'production', 
      sameSite: 'Strict', 
    });
    sessionStorage.setItem('user', response.data.data);

    return response.data.data; 
  } catch (error) {
    console.error('Lỗi khi đăng nhập:', error);
    throw error; 
  }
};
export const register = async (username, email, phone, password) => {
  return axios.post('/api/register', { username, email, phone, password });
};