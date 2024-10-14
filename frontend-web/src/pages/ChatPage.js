import React, { useState, useEffect } from 'react';
import Sidebar from '../components/Messages/MessSidebar/SidebarChat';
import ChatContent from '../components/Messages/ChatContent/ChatContent';
import Layout from '../layouts/Layout';
import '../assets/css/ChatPage.css'; 

function ChatPage() {
  const [userData, setUserData] = useState(null); 
  const [messages, setMessages] = useState([
    { name: 'Nguyễn Huy Hoàng', avatar: 'https://firebasestorage.googleapis.com/v0/b/multi-aura.appspot.com/o/Hihon%2F1728534046_9ea1c9841cadbef3e7bc.jpg?alt=media&token=3d221a08-d064-4ece-881a-32e2c5d273e1', lastMessage: 'Nguyễn sent an attachment · 1w' },
    { name: 'Vương Kim Đinh', avatar: 'https://firebasestorage.googleapis.com/v0/b/multi-aura.appspot.com/o/Hihon%2F1728534046_9ea1c9841cadbef3e7bc.jpg?alt=media&token=3d221a08-d064-4ece-881a-32e2c5d273e1', lastMessage: 'Đauuuu chứ · 33m' },
  ]);

  const [currentChat, setCurrentChat] = useState({
    user: { name: 'Vương Kim Đinh', avatar: 'https://firebasestorage.googleapis.com/v0/b/multi-aura.appspot.com/o/Hihon%2F1728534046_9ea1c9841cadbef3e7bc.jpg?alt=media&token=3d221a08-d064-4ece-881a-32e2c5d273e1', status: 'Active 34m ago' },
    messages: [
      { text: 'Đauuuu chứ :)))', time: '23 Jul 2024, 02:41', isSentByUser: false },
      { text: 'Đi mạnh giỏi nhaaaa', time: '23 Jul 2024, 06:07', isSentByUser: true },
    ]
  });

  useEffect(() => {
    const storedUser = localStorage.getItem('user'); 
    if (storedUser) {
      setUserData(JSON.parse(storedUser)); 
    }
  }, []);

  const handleSelectChat = (message) => {
    const selectedUserChat = {
      user: { name: message.name, avatar: message.avatar, status: 'Active recently' },
      messages: [
        { text: 'This is a previous message', time: '22 Jul 2024, 04:41', isSentByUser: false },
        { text: 'Another old message', time: '22 Jul 2024, 06:07', isSentByUser: true },
        { text: 'New message here', time: '23 Jul 2024, 09:07', isSentByUser: true }
      ]
    };
    setCurrentChat(selectedUserChat);
  };

  return (
    <Layout userData={userData}>
      <div className="container-fluid chat-page">
        <div className="row">
          <div className="col-lg-3 col-md-3 col-sm-12 sidebar-wrapper"> 
            <Sidebar messages={messages} onSelectChat={handleSelectChat} />
          </div>
          <div className="col-lg-9 col-md-9 col-sm-12 chat-content-wrapper">
          <ChatContent chat={currentChat} />
          </div>
        </div>
      </div>
    </Layout>
  );
}

export default ChatPage;
