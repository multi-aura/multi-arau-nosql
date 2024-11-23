import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faBellSlash,
  faThumbtack,
  faUsers,
  faClock,
  faFile,
  faImage,
  faTimes
} from '@fortawesome/free-solid-svg-icons';
import './SettingSidebarChat.css';

const SettingSidebarChat = ({ isOpen, onClose, currentChat }) => {
  return (
    <div className={`setting-sidebar-chat ${isOpen ? 'visible' : 'hidden'}`}>
      {/* Header */}
      <div className="header">
        <img
          src={currentChat?.avatar || 'default-avatar.png'}
          alt="Avatar"
        />
        <h3>{currentChat?.name_conversation || 'Hội thoại'}</h3>
        <div className="actions">
          <button>
            <FontAwesomeIcon icon={faBellSlash} /> Tắt thông báo
          </button>
          <button>
            <FontAwesomeIcon icon={faThumbtack} /> Ghim hội thoại
          </button>
          <button>
            <FontAwesomeIcon icon={faUsers} /> Tạo nhóm
          </button>
        </div>
      </div>

      {/* Danh sách nhắc nhở */}
      <div className="section">
        <h4 className="section-title">
          <FontAwesomeIcon icon={faClock} /> Danh sách nhắc nhở
        </h4>
        <p>Không có nhắc nhở nào</p>
      </div>

      {/* Ảnh/Video */}
      <div className="section">
        <h4 className="section-title">
          <FontAwesomeIcon icon={faImage} /> Ảnh/Video
        </h4>
        <div className="media-grid">
          {/* Thêm ảnh */}
          <img src="image1.jpg" alt="Media 1" />
          <img src="image2.jpg" alt="Media 2" />
          <img src="image3.jpg" alt="Media 3" />
          <img src="image4.jpg" alt="Media 4" />
        </div>
        <div className="view-all-btn">Xem tất cả</div>
      </div>

      {/* File */}
      <div className="section">
        <h4 className="section-title">
          <FontAwesomeIcon icon={faFile} /> File
        </h4>
        <ul className="file-list">
          <li>
            <div className="file-name">
              <FontAwesomeIcon icon={faFile} />
              <span>2001216069_File1.xlsx</span>
            </div>
            <div className="file-date">21/11/2024</div>
          </li>
          <li>
            <div className="file-name">
              <FontAwesomeIcon icon={faFile} />
              <span>Nhom10_BaoCaoCuoiKi.zip</span>
            </div>
            <div className="file-date">11/11/2024</div>
          </li>
          <li>
            <div className="file-name">
              <FontAwesomeIcon icon={faFile} />
              <span>EMHUNer.docx</span>
            </div>
            <div className="file-date">10/11/2024</div>
          </li>
        </ul>
        <div className="view-all-btn">Xem tất cả</div>
      </div>

     
    </div>
  );
};

export default SettingSidebarChat;
