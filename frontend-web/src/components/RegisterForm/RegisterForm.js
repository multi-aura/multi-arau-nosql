import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faUser, faLock, faEnvelope, faPhone } from '@fortawesome/free-solid-svg-icons';
import CustomInput from '../Input/CustomInput'; // Giữ nguyên component CustomInput từ login
import { validateEmail } from '../../utils/validation'; // Hàm validateEmail
import './RegisterForm.css';

function RegisterForm({ handleRegister }) {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [phone, setPhone] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [errorMessage, setErrorMessage] = useState('');
  const [birthday, setErrorbirthday] = useState('')
  const onSubmit = (e) => {
    e.preventDefault();

    if (password !== confirmPassword) {
      setErrorMessage('Mật khẩu và mật khẩu xác nhận không khớp.');
      return;
    }

    if (!validateEmail(email)) {
      setErrorMessage('Email không hợp lệ.');
      return;
    }

    handleRegister(username, email, phone, password);
  };

  return (
    <form onSubmit={onSubmit} className="custom-register-form">
      <h2 className="form-title text-center mb-4">Register</h2>

      {errorMessage && <p className="text-danger text-center">{errorMessage}</p>} {/* Hiển thị thông báo lỗi */}

      <CustomInput
        type="text"
        label="Username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        icon={faUser}
      />
      <CustomInput
        type="email"
        label="Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        icon={faEnvelope}
      />
      <CustomInput
        type="text"
        label="Phone"
        value={phone}
        onChange={(e) => setPhone(e.target.value)}
        icon={faPhone}
      />
      <CustomInput
        type="password"
        label="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        icon={faLock}
      />
      <CustomInput
        type="password"
        label="Confirm Password"
        value={confirmPassword}
        onChange={(e) => setConfirmPassword(e.target.value)}
        icon={faLock}
      />

      <button type="submit" className="custom-button w-100">Register</button>
    </form>
  );
}

export default RegisterForm;
