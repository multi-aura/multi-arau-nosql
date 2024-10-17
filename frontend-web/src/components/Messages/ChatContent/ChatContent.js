import React, { useState, useEffect } from 'react';
import ChatHeader from '../ChatHeader/ChatHeader';
import MessageBubble from '../MessageBubble/MessageBubble';
import ChatInput from '../ChatInput/ChatInput';
import './ChatContent.css';

const ChatContent = ({ chat, onSendMessage }) => {
  const [messages, setMessages] = useState([]);

  // Sắp xếp và cập nhật tin nhắn mỗi khi `chat` thay đổi
  useEffect(() => {
    if (chat.chats && Array.isArray(chat.chats)) {
      const sortedMessages = [...chat.chats].sort((a, b) => new Date(a.createdat) - new Date(b.createdat));
      setMessages(sortedMessages);
    }
  }, [chat]);

  const handleSendMessage = async (messageContent) => {
    try {
      // Gửi tin nhắn qua API và đợi phản hồi từ API (sau khi tin nhắn đã lưu vào DB)
      const newMessage = await onSendMessage(messageContent);

      // Chỉ khi API trả về thành công, bạn mới cập nhật giao diện
      if (newMessage) {
        setMessages(prevMessages => [...prevMessages, newMessage]); // Cập nhật danh sách tin nhắn với tin nhắn mới
      }
    } catch (error) {
      // Nếu xảy ra lỗi, hiển thị thông báo lỗi hoặc log ra console
      console.error('Lỗi khi gửi tin nhắn:', error);
      alert('Có lỗi xảy ra khi gửi tin nhắn. Vui lòng thử lại.');
    }
  };


  return (
    <div className="chat-content">
      <ChatHeader user={chat} />

      <div className="chat-messages">
        {messages.length > 0 ? (
          messages.map((message, index) => (
            <MessageBubble
              key={index}
              message={message}
              userAvatar={message.sender?.avatar || 'default-avatar.png'}   // Đặt ảnh mặc định nếu sender không tồn tại
            />
          ))
        ) : (
          <div className="no-message-container">
            <p className="no-message-text">
              Chưa có tin nhắn nào... nhưng đây chỉ là sự khởi đầu! Hãy gửi một lời chào thật ấm áp hoặc một câu chuyện thú vị để bắt đầu cuộc trò chuyện.
            </p>
          </div>
        )}
      </div>

      {/* Phần nhập tin nhắn */}
      <ChatInput onSendMessage={handleSendMessage} />
    </div>
  );
};

export default ChatContent;
