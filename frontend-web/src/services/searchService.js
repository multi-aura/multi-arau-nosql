import axios from 'axios';
import { API_URL } from '../config/config';
import Cookies from 'js-cookie';

const SEARCH_URL = `${API_URL}/search`;
const SEARCH_PEOPLE_URL = `${SEARCH_URL}/people`;

export const searchPeople = async (query, limit = 4, page = 1) => {
    try {
        const token = Cookies.get('authToken');
        const response = await axios.post(`${SEARCH_PEOPLE_URL}?q=${query}`, 
            {
                limit: limit,
                page: page
            },
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
        return response.data;
    } catch (error) {
        console.error('Fail to search:', error);
        throw error;
    }
};

export const getPeopleSuggestions = async (limit = 6, page = 1)  => {
    try {
        const token = Cookies.get('authToken');
        if (!token) {
            throw new Error('Token không tồn tại. Vui lòng đăng nhập.');
        }

        const response = await axios.post(SEARCH_PEOPLE_URL, 
            {
                limit: limit,
                page: page,
            },
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
    

        return response.data;
    } catch (error) {
        console.error('Lỗi khi lấy danh sách gợi ý:', error.response ? error.response.data : error.message);
        throw error;
    }
};

export const getForYouPosts = async (limit = 10, page = 1) => {
    try {
      const token = Cookies.get('authToken'); 
  
      const response = await axios.post(`${SEARCH_URL}/for-you`, 
        {
          limit, 
          page
        },
        {
          headers: {
            Authorization: `Bearer ${token}` 
          }
        }
      );
      
      return response.data;
    } catch (error) {
      console.error('Lỗi khi lấy bài viết "For You":', error);
      throw error;
    }
  };

  export const getNewsPosts = async (limit = 10, page = 1) => {
    try {
      const token = Cookies.get('authToken'); 
  
      const response = await axios.post(`${SEARCH_URL}/news`, 
        {
          limit, 
          page
        },
        {
          headers: {
            Authorization: `Bearer ${token}` 
          }
        }
      );
      
      return response.data;
    } catch (error) {
      console.error('Lỗi khi lấy bài viết "News":', error);
      throw error;
    }
  };
  export const getTrendingPosts = async (limit = 10, page = 1) => {
    try {
      const token = Cookies.get('authToken'); 
  
      const response = await axios.post(`${SEARCH_URL}/trending`, 
        {
          limit, 
          page
        },
        {
          headers: {
            Authorization: `Bearer ${token}` 
          }
        }
      );
      
      return response.data;
    } catch (error) {
      console.error('Lỗi khi lấy bài viết "Trending":', error);
      throw error;
    }
  };

  export const getPosts = async (limit = 10, page = 1) => {
    try {
      const token = Cookies.get('authToken'); 
  
      const response = await axios.post(`${SEARCH_URL}/posts`, 
        {
          limit, 
          page
        },
        {
          headers: {
            Authorization: `Bearer ${token}` 
          }
        }
      );
      
      return response.data;
    } catch (error) {
      console.error('Lỗi khi lấy bài viết "Post":', error);
      throw error;
    }
  };