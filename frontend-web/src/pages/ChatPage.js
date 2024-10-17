import React, { useState, useEffect, useRef } from 'react';  // Thêm useRef vào đây

import Sidebar from '../components/Messages/MessSidebar/SidebarChat';
import ChatContent from '../components/Messages/ChatContent/ChatContent';
import Layout from '../layouts/Layout';
import '../assets/css/ChatPage.css';
import { getUserConversation, getConversationDetails, sendMessageToConversation } from "../services/chatservice";
import { MdMargin } from 'react-icons/md';

function ChatPage() {
  const [userData, setUserData] = useState(null);
  const [messages, setMessages] = useState([]);
  const [currentChat, setCurrentChat] = useState(null);
  const [loadingChat, setLoadingChat] = useState(false);
  const ws = useRef(null); // Tạo WebSocket ref để giữ kết nối
  useEffect(() => {
    ws.current = new WebSocket('ws://localhost:3000/ws');


    ws.current.onopen = () => {
      console.log('WebSocket connection established');
    };

    ws.current.onmessage = (event) => {
      const newMessage = JSON.parse(event.data);
      if (newMessage.conversationID === currentChat?._id) {
        // Cập nhật tin nhắn mới vào giao diện
        setMessages((prevMessages) => [...prevMessages, newMessage]);
      }
    };

    ws.current.onclose = () => {
      console.log('WebSocket connection closed');
    };

    ws.current.onerror = (error) => {
      console.log('WebSocket error:', error);
    };

    // Đóng kết nối WebSocket khi component unmount
    return () => {
      if (ws.current) {
        ws.current.close();
      }
    };
  }, [currentChat]); // Reconnect WebSocket when currentChat changes

  // Lấy thông tin người dùng từ localStorage
  useEffect(() => {
    const storedUser = localStorage.getItem('user');
    if (storedUser) {
      setUserData(JSON.parse(storedUser));
    }
  }, []);

  // Lấy danh sách các cuộc trò chuyện của người dùng
  useEffect(() => {
    const fetchUserConversation = async () => {
      try {
        if (userData) {
          const userConversations = await getUserConversation(userData.userID);
          setMessages(userConversations);
        }
      } catch (error) {
        console.error('Error fetching conversation:', error);
      }
    };
    if (userData) {
      fetchUserConversation();
    }
  }, [userData]);

  // Hàm khi người dùng chọn một cuộc trò chuyện từ Sidebar
  const handleSelectChatMessage = async (conversationID) => {
    setLoadingChat(true);
    try {

      const conversationData = await getConversationDetails(conversationID);
      if (conversationData.users && conversationData.users.length > 0) {

        setCurrentChat(conversationData);
      } else {
        console.error('No users found in the conversation');
        setCurrentChat(false);
      }

      setLoadingChat(false);
    } catch (error) {
      console.error('Error fetching conversation details:', error);
      setLoadingChat(false);
      setCurrentChat(false);
    }
  };
  // Hàm gửi tin nhắn
  const handleSendMessage = async (messagesContent) => {
    if (!currentChat || !userData || !ws.current) return;

    // Chuẩn bị dữ liệu tin nhắn
    let content = {
      text: messagesContent.text || "",
      image: messagesContent.image || "",
      voice_url: messagesContent.voice_url || ""
    };

    const messageData = {
      conversationID: currentChat._id,
      user_id: userData.userID,
      content: content
    };

    try {
      // Gửi tin nhắn qua API để lưu vào DB
      const savedMessage = await sendMessageToConversation(currentChat._id, userData.userID, content);

      // Kiểm tra trạng thái WebSocket trước khi gửi tin nhắn
      if (ws.current.readyState === WebSocket.OPEN) {
        // Gửi tin nhắn qua WebSocket tới tất cả các client khác
        ws.current.send(JSON.stringify(savedMessage));
      } else {
        console.error("WebSocket connection is not open.");
      }

      // Cập nhật tin nhắn ngay lập tức cho người dùng hiện tại
      setMessages((prevMessages) => [...prevMessages, savedMessage]);

    } catch (error) {
      console.error("Error sending message:", error);
    }
  };



  return (
    <Layout userData={userData}>
      <div className="container-fluid chat-page">
        <div className="row">
          <div className="col-lg-3 col-md-3 col-sm-12 sidebar-wrapper">
            {/* Truyền handleSelectChatMessage thay cho handleSelectChat */}
            <Sidebar messages={messages} onSelectChat={handleSelectChatMessage} />
          </div>
          <div className="col-lg-9 col-md-9 col-sm-12 chat-content-wrapper" style={{ height: "90vh", padding: '0px 10px' }}>
            {loadingChat ? (
              <div>Loading chat...</div>
            ) : currentChat ? (
              <ChatContent chat={currentChat} onSendMessage={handleSendMessage} />
            ) : (
              <div>Please select a chat to view</div>
            )}
          </div>
        </div>
      </div>
    </Layout>
  );
}

export default ChatPage;
