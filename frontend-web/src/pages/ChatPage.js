import React, { useState, useEffect, useRef } from 'react';
import Sidebar from '../components/Messages/MessSidebar/SidebarChat';
import ChatContent from '../components/Messages/ChatContent/ChatContent';
import Layout from '../layouts/Layout';
import '../assets/css/ChatPage.css';
import { getUserConversation, getConversationDetails, sendMessageToConversation } from "../services/chatservice";

function ChatPage() {
  const [userData, setUserData] = useState(null);
  const [conversations, setConversations] = useState([]);  // Tạo state riêng cho danh sách các cuộc trò chuyện
  const [messages, setMessages] = useState([]);  // Chỉ dành cho tin nhắn

  const [currentChat, setCurrentChat] = useState(null);
  const [loadingChat, setLoadingChat] = useState(false);
  const ws = useRef(null);
  const [newMessageItems, setNewMessageItems] = useState(null);
  useEffect(() => {
    const storedUser = localStorage.getItem('user');
    if (storedUser) {
      setUserData(JSON.parse(storedUser));
    }
  }, []);

  // Kết nối WebSocket khi component được mount
  useEffect(() => {
    if (userData && currentChat) {
      ws.current = new WebSocket(`ws://localhost:3002/ws?user_id=${userData.userID}&conversation_id=${currentChat._id}`);

      // Lắng nghe sự kiện khi có tin nhắn mới từ server qua WebSocket
      ws.current.onmessage = (event) => {
        const receivedMessage = JSON.parse(event.data);

        // Cập nhật messages với tin nhắn mới
        if (receivedMessage.conversationID === currentChat._id) {
          setMessages((prevMessages) => {
            // Đảm bảo prevMessages luôn là một mảng hợp lệ
            if (!Array.isArray(prevMessages)) {
              prevMessages = [];  // Đặt lại thành mảng rỗng nếu không phải là mảng
            }

            const updatedMessages = [...prevMessages, receivedMessage];

            return updatedMessages;
          });
          setNewMessageItems(receivedMessage);


        } else {
          console.log("Message does not belong to the current conversation.");
        }
      };

      // Xử lý sự kiện khi kết nối bị đóng
      ws.current.onclose = () => {
        console.log("WebSocket disconnected");
      };

      return () => {
        if (ws.current) ws.current.close();
      };
    }
  }, [userData, currentChat]);

  // Lấy danh sách các cuộc trò chuyện của người dùng
  useEffect(() => {
    const fetchUserConversation = async () => {
      try {
        if (userData) {
          const userConversations = await getUserConversation(userData.userID);
          setConversations(userConversations);  // Lưu các cuộc trò chuyện vào `conversations`
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
        setMessages(conversationData.chats || []);  // Cập nhật tin nhắn của cuộc trò chuyện
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

  const handleSendMessage = async (messageContent) => {
    // Cấu hình nội dung tin nhắn
    let content = {
      text: messageContent.text || "",
      image: messageContent.image || "",
      voice_url: messageContent.voice_url || ""
    };

    const messageData = {
      conversationID: currentChat._id,
      sender: {
        userID: userData.userID,
        fullname: userData.fullname || "Unknown",
        username: userData.username || "",
        avatar: userData.avatar || "default-avatar-url",  // Link avatar của người gửi
        added_at: new Date().toISOString()  // Thời gian hiện tại
      },
      content: content,  // Nội dung của tin nhắn
      createdat: new Date().toISOString(),  // Thời gian tin nhắn được tạo
      updatedat: new Date().toISOString(),  // Thời gian tin nhắn được cập nhật
      status: "sent"  // Trạng thái của tin nhắn
    };
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(messageData));  // Gửi tin nhắn qua WebSocket
    } else {
      console.error("WebSocket is not open.");
    }

    try {
      await sendMessageToConversation(currentChat._id, userData.userID, content);
    } catch (error) {
      console.error("Error sending message:", error);
    }
    // Cập nhật tin nhắn mới và truyền qua Sidebar
    setNewMessageItems(messageData);
  };



  return (
    <Layout userData={userData}>
      <div className="container-fluid chat-page">
        <div className="row">
          <div className="col-lg-3 col-md-3 col-sm-12 sidebar-wrapper">
            <Sidebar conversations={conversations} onSelectChat={handleSelectChatMessage} newMessageItems={newMessageItems} />
          </div>
          <div className="col-lg-9 col-md-9 col-sm-12 chat-content-wrapper" style={{ height: "90vh", padding: '0px 0px' }}>
            {loadingChat ? (
              <div>Loading chat...</div>
            ) : currentChat ? (
              <>
                <ChatContent chat={currentChat} messages={messages} currentUserID={userData.userID} onSendMessage={handleSendMessage} />
              </>
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
