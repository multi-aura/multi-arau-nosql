import React, { useState } from 'react';

// Hàm kiểm tra URL hình ảnh hợp lệ
const isValidImageURL = (url) => {
  return /\.(jpg|jpeg|png|gif)$/i.test(url);
};

// Hàm kiểm tra URL âm thanh hợp lệ
const isValidAudioURL = (url) => {
  return /\.(mp3|wav)$/i.test(url);
};

const ChatInput = ({ onSendMessage }) => {
  const [message, setMessage] = useState('');

  const handleSendMessage = () => {
    if (message.trim()) {
      let messageContent = {};

      // Kiểm tra nếu là hình ảnh
      if (isValidImageURL(message.trim())) {
        messageContent = { image: message.trim() }; // Nếu là URL hình ảnh, gán vào image
      }
      // Kiểm tra nếu là âm thanh
      else if (isValidAudioURL(message.trim())) {
        messageContent = { voice_url: message.trim() }; // Nếu là URL âm thanh, gán vào voice_url
      }
      // Ngược lại, coi như là text
      else {
        messageContent = { text: message.trim() }; // Nếu là văn bản, gán vào text
      }

      // Gửi nội dung tin nhắn (văn bản, hình ảnh hoặc âm thanh)
      onSendMessage(messageContent);
      setMessage(''); // Xóa nội dung input sau khi gửi
    }
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleSendMessage();
    }
  };

  return (
    <div className="chat-input">
      <input
        type="text"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder="Nhập tin nhắn, URL hình ảnh hoặc âm thanh..."
        className="input-message"
      />
      <button onClick={handleSendMessage} disabled={!message.trim()}>
        Gửi
      </button>
    </div>
  );
};

export default ChatInput;
