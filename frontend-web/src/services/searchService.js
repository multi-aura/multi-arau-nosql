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
