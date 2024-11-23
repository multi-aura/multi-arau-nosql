import React, { useState, useEffect, useRef } from 'react';
import Sidebar from '../components/Messages/MessSidebar/SidebarChat';
import ChatContent from '../components/Messages/ChatContent/ChatContent';
import SettingSidebarChat from '../components/Messages/SettingSidebarChat/SettingSidebarChat';
import Layout from '../layouts/Layout';
import '../assets/css/ChatPage.css';
import { getUserConversation, getConversationDetails, sendMessageToConversation } from "../services/chatservice";
import { API_URL_WS } from '../config/config';

function ChatPage() {
  const [userData, setUserData] = useState(null);
  const [conversations, setConversations] = useState([]);
  const [messages, setMessages] = useState([]);
  const [currentChat, setCurrentChat] = useState(null);
  const [loadingChat, setLoadingChat] = useState(false);
  const [isSidebarOpen, setIsSidebarOpen] = useState(false); // Trạng thái mở/đóng sidebar
  const ws = useRef(null);
  const [newMessageItems, setNewMessageItems] = useState(null);

  // Load dữ liệu người dùng từ localStorage
  useEffect(() => {
    const storedUser = localStorage.getItem('user');
    if (storedUser) {
      setUserData(JSON.parse(storedUser));
    }
  }, []);

  // Kết nối WebSocket khi có userData và currentChat
  useEffect(() => {
    if (userData && currentChat) {
      ws.current = new WebSocket(`${API_URL_WS}/ws?user_id=${userData.userID}&conversation_id=${currentChat._id}`);

      ws.current.onmessage = (event) => {
        const receivedMessage = JSON.parse(event.data);
        if (receivedMessage.conversationID === currentChat._id) {
          setMessages((prevMessages) => [...prevMessages, receivedMessage]);
          setNewMessageItems(receivedMessage);
        } else {
          console.log("Message does not belong to the current conversation.");
        }
      };

      ws.current.onclose = () => {
        console.log("WebSocket disconnected");
      };

      return () => {
        if (ws.current) ws.current.close();
      };
    }
  }, [userData, currentChat]);

  // Lấy danh sách cuộc trò chuyện của người dùng
  useEffect(() => {
    const fetchUserConversation = async () => {
      try {
        if (userData) {
          const userConversations = await getUserConversation(userData.userID);
          setConversations(userConversations);
        }
      } catch (error) {
        console.error('Error fetching conversation:', error);
      }
    };
    if (userData) {
      fetchUserConversation();
    }
  }, [userData]);

  const handleSelectChatMessage = async (conversationID) => {
    setLoadingChat(true);
    try {
      const conversationData = await getConversationDetails(conversationID, userData.userID);
      setCurrentChat(conversationData);
      setMessages(conversationData.chats || []);
      setLoadingChat(false);
    } catch (error) {
      console.error('Error fetching conversation details:', error);
      setLoadingChat(false);
    }
  };

  const handleSendMessage = async (messageContent) => {
    const content = {
      text: messageContent.text || "",
      image: messageContent.image || "",
      voice_url: messageContent.voice_url || ""
    };

    const messageData = {
      conversationID: currentChat._id,
      sender: {
        userID: userData.userID,
        fullname: userData.fullname || "Unknown",
        avatar: userData.avatar || "default-avatar-url"
      },
      content,
      createdat: new Date().toISOString()
    };

    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(messageData));
    } else {
      console.error("WebSocket is not open.");
    }

    try {
      await sendMessageToConversation(currentChat._id, userData.userID, content);
    } catch (error) {
      console.error("Error sending message:", error);
    }
    setNewMessageItems(messageData);
  };

  return (
    <Layout userData={userData}>
      <div className={`container-fluid chat-page ${isSidebarOpen ? 'with-sidebar' : ''}`}>
        <div className="row">
          <div className="col-lg-3 col-md-3 col-sm-12 sidebar-wrapper">
            <Sidebar
              conversations={conversations}
              onSelectChat={handleSelectChatMessage}
              newMessageItems={newMessageItems}
            />
          </div>
          <div
            className={`col-lg-9 col-md-9 col-sm-12 chat-content-wrapper`}
            style={{
              display: 'flex',
              height: '90vh',
              padding: '0',
            }}
          >
            <div
              className={`chat-content ${isSidebarOpen ? 'shrink-content' : ''}`}
              style={{ flex: isSidebarOpen ? '0.7' : '1', transition: 'flex 0.3s ease' }}
            >
              {loadingChat ? (
                <div>Loading chat...</div>
              ) : currentChat ? (
                <ChatContent
                chat={currentChat} // Truyền currentChat xuống
                messages={messages}
                currentUserID={userData.userID}
                onSendMessage={handleSendMessage}
                onToggleSidebar={() => setIsSidebarOpen(!isSidebarOpen)} // Truyền hàm mở/đóng sidebar
                isSidebarOpen={isSidebarOpen} // Trạng thái sidebar
                userData={userData} // Truyền userData xuống
              />
              
              ) : (
                <div>Please select a chat to view</div>
              )}
            </div>

            {/* Sidebar */}
            {isSidebarOpen && (
              <SettingSidebarChat
                isOpen={isSidebarOpen}
                onClose={() => setIsSidebarOpen(false)}
                currentChat={currentChat}
              />
            )}
          </div>
        </div>
      </div>
    </Layout>
  );
}

export default ChatPage;
