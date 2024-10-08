import React from 'react';
import { useNavigate } from 'react-router-dom';
import '../assets/css/IntroPage.css'; // CSS tùy chỉnh cho trang Intro

function IntroPage() {
  const navigate = useNavigate();

  const handleStart = () => {
    navigate('/login'); // Chuyển hướng đến trang đăng nhập
  };

  return (
    <div className="intro-page-container" onClick={handleStart}>
      <div className="intro-content">
        <h1>Multi Aura</h1>
        <p>Start Joining Us</p>
      </div>
    </div>
  );
}

export default IntroPage;
