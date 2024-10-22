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

export const getConversationDetails = async (conversationID) => {
    try {
        const token = Cookies.get('authToken');
        const response = await axios.get(`${CONVERSATION_URL}/${conversationID}`, {

            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        return response.data.data;
    } catch (error) {
        console.log(`Failed to get conversation details for conversation with ID ${conversationID}`, error);
        throw error;
    }
};
export const sendMessageToConversation = async (conversationID, userID, content) => {

    try {
        const token = Cookies.get('authToken');
        const response = await axios.post(`${CONVERSATION_URL}/send-message/${conversationID}`, {
            user_id: userID,
            content: content
        }, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });


        // Kiểm tra nếu response.data.data không tồn tại, bạn có thể log một cảnh báo
        if (!response.data.data) {
            console.warn('API response does not contain expected "data" field:', response.data);
            return null;
        }

        return response.data.data; // Trả về dữ liệu tin nhắn
    } catch (error) {
        console.log(`Failed to send message for conversation with ID ${conversationID}`, error);
        throw error;
    }
};
