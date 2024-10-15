import axios from 'axios';
import { API_URL } from '../config/config';
import Cookies from 'js-cookie';

const CONVERSATION_URL = `${API_URL}/conversation`;

export const getUserConversation = async (userID) => {
    try {
        const token = Cookies.get('authToken');
        const response = await axios.get(`${CONVERSATION_URL}/get-user-conversations/${userID}`, {
            headers: {
                Authorization: `Bearer ${token}` 
            }
        });
        return response.data.data;

    } catch (error) {
        console.log(`Failed to get conversation for user with ID ${userID}`, error);
        throw error;
    }
}