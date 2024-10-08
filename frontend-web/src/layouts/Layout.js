import React from 'react';
import Sidebar from '../components/Sidebar/Sidebar';
import Header from '../components/Header/Header';
import './Layout.css';
function Layout({ children, userData }) {
  return (
    <div className="container-fluid p-0 w-full">
      <div className="row no-gutters">
        <div className="col-2 bg-dark text-white p-0 w-full">
          <Sidebar />
        </div>

        <div className="col-10 bg-dark d-flex flex-column p-0">
          <Header userData={userData} />
          <div className="content-area bg-dark flex-grow-1 p-1">
            {children}
          </div>
        </div>
      </div>
    </div>
  );
}

export default Layout;
