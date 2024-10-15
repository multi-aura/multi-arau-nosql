import React, { useState, useEffect } from 'react';
import SearchBar from '../components/Explore/Search/SearchExploreBar';
import TabMenu from '../components/Explore/TabMenu/TabMenu';
import SuggestedUsers from '../components/Explore/SuggestedUser/ListSuggestedUsers';
import Layout from '../layouts/Layout';
import PeopleSearchResult from '../components/Explore/ExploreSubPage/PeopleSearchResult';
import '../assets/css/Explore.css';

function Explore() {
  const [userData, setUserData] = useState(null); 
  const [suggestedUsers, setSuggestedUsers] = useState([]);
  const [activeTab, setActiveTab] = useState('Trending');
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    const storedUser = localStorage.getItem('user'); 
    if (storedUser) {
      setUserData(JSON.parse(storedUser)); 
    }

    const suggested = [
      { id: 1, name: 'Nguyễn Huy Hoàng', avatar: 'https://firebasestorage.googleapis.com/v0/b/multi-aura.appspot.com/o/Hihon%2F1728534046_9ea1c9841cadbef3e7bc.jpg?alt=media&token=3d221a08-d064-4ece-881a-32e2c5d273e1' },
      { id: 2, name: 'Kim Đinh', avatar: 'https://firebasestorage.googleapis.com/v0/b/multi-aura.appspot.com/o/Hihon%2F1728534046_9ea1c9841cadbef3e7bc.jpg?alt=media&token=3d221a08-d064-4ece-881a-32e2c5d273e1' }
    ];
    setSuggestedUsers(suggested);
  }, []); 

  const renderContent = () => {
    switch (activeTab) {
      case 'For you':
        return <div><h2>Nội dung "For You" sẽ hiển thị ở đây</h2></div>;
      case 'Trending':
        return <div><h2>Nội dung "Trending" sẽ hiển thị ở đây</h2></div>;
      case 'News':
        return <div><h2>Nội dung "News" sẽ hiển thị ở đây</h2></div>;
      case 'People':
        return <PeopleSearchResult suggestedUsers={suggestedUsers} />;
      case 'Posts':
        return <div><h2>Nội dung "Posts" sẽ hiển thị ở đây</h2></div>;
      default:
        return <div><h2>Nội dung mặc định sẽ hiển thị ở đây</h2></div>;
    }
  };
  return (
    <Layout userData={userData}>
      <div className="explore-page container-fluid">
        <div className="row">
          <div className="col-lg-8 col-md-7 col-sm-12 mb-4">
            <SearchBar onSearch={setSearchTerm} />
            <TabMenu activeTab={activeTab} setActiveTab={setActiveTab} />
            <div className="post-container">
              {renderContent()}
            </div>
          </div>

          <div className="col-lg-4 col-md-5 col-sm-12 ">
            <div className="suggestions-container p-3 rounded">
              <SuggestedUsers suggestedUsers={suggestedUsers} />
            </div>
          </div>
        </div>
      </div>
    </Layout>
  );
}

export default Explore;
