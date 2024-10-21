import React from 'react';
import ChatHeader from '../ChatHeader/ChatHeader';
import MessageBubble from '../MessageBubble/MessageBubble';
import ChatInput from '../ChatInput/ChatInput';
import './ChatContent.css';

const ChatContent = ({ chat, messages, onSendMessage, currentUserID }) => {
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
      <ChatHeader user={chat} />
      <div className="chat-messages">
        {messages.length > 0 ? (
          messages.map((message, index) => {
            // Kiểm tra tin nhắn trước và tin nhắn tiếp theo
            const previousMessage = messages[index - 1];
            const nextMessage = messages[index + 1];
            const isSameSenderAsPrevious = previousMessage && previousMessage.sender.userID === message.sender.userID;
            const isLastMessageFromSameSender = !nextMessage || nextMessage.sender.userID !== message.sender.userID;

            // Hiển thị tên người gửi nếu là tin nhắn đầu tiên trong chuỗi của họ
            const showSenderInfo = !isSameSenderAsPrevious;

            return (
              <MessageBubble
                key={index}
                message={message}
                userAvatar={message.sender?.avatar || 'default-avatar.png'}
                currentUserID={currentUserID}
                showSenderInfo={showSenderInfo} // Hiển thị tên người gửi khi cần
                showTime={isLastMessageFromSameSender} // Hiển thị thời gian cho tin nhắn cuối cùng trong chuỗi
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
      </div>
      <ChatInput onSendMessage={handleSendMessage} />
    </div>
  );
};

export default ChatContent;
