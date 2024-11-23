import React, { useEffect, useRef } from 'react';
import ChatHeader from '../ChatHeader/ChatHeader';
import MessageBubble from '../MessageBubble/MessageBubble';
import ChatInput from '../ChatInput/ChatInput';
import './ChatContent.css';
const ChatContent = ({
  chat, // Prop chứa thông tin chat hiện tại
  messages,
  onSendMessage,
  currentUserID,
  onToggleSidebar, // Hàm bật/tắt sidebar
  isSidebarOpen, // Trạng thái sidebar
  userData // Dữ liệu người dùng
}) => {
  const messageEndRef = useRef(null); // Tạo tham chiếu để cuộn đến cuối

  // Tự động cuộn xuống cuối khi có tin nhắn mới
  useEffect(() => {
    messageEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSendMessage = async (messageContent) => {
    try {
      const newMessage = await onSendMessage(messageContent);
      if (newMessage) {
        // Sau khi gửi tin nhắn thành công, parent component sẽ tự cập nhật `messages`
      }
    } catch (error) {
      console.error('Error sending message:', error);
      alert('Có lỗi xảy ra khi gửi tin nhắn. Vui lòng thử lại.');
    }
  };

  return (
    <div className="chat-content">
      <ChatHeader
        user={chat}
        currentUserID={userData.userID}
        onToggleSidebar={onToggleSidebar} // Truyền hàm bật/tắt sidebar
        isSidebarOpen={isSidebarOpen} // Trạng thái sidebar
      />

      <div className="chat-messages">
        {messages.length > 0 ? (
          messages.map((message, index) => {
            const previousMessage = messages[index - 1];
            const nextMessage = messages[index + 1];
            const isSameSenderAsPrevious = previousMessage && previousMessage.sender.userID === message.sender.userID;
            const isLastMessageFromSameSender = !nextMessage || nextMessage.sender.userID !== message.sender.userID;
            const showSenderInfo = !isSameSenderAsPrevious;

            return (
              <MessageBubble
                key={index}
                message={message}
                userAvatar={message.sender?.avatar || 'default-avatar.png'}
                currentUserID={currentUserID}
                showSenderInfo={showSenderInfo}
                showTime={isLastMessageFromSameSender}
              />
            );
          })
        ) : (
          <div className="no-message-container">
            <p className="no-message-text">
              Chưa có tin nhắn nào... nhưng đây chỉ là sự khởi đầu! Hãy gửi một lời chào thật ấm áp hoặc một câu chuyện thú vị để bắt đầu cuộc trò chuyện.
            </p>
          </div>
        )}
        <div ref={messageEndRef} />
      </div>
      <ChatInput onSendMessage={handleSendMessage} />
    </div>
  );
};

export default ChatContent;
