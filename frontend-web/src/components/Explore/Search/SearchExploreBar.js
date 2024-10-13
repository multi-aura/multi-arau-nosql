import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSearch } from '@fortawesome/free-solid-svg-icons';
import './SearchExploreBar.css';
import SearchResults from './SearchResults';

const SearchBar = () => {
  const [searchTerm, setSearchTerm] = useState('');
  const [showResults, setShowResults] = useState(false);

  const recentSearches = ['Nguyễn Huy Hoàng', 'Kim Đinh', 'Minh Thư']; // Ví dụ Recent
  const suggestions = ['Nguyễn Huy Hoàng', 'Kim Đinh', 'Minh Thư']; // Ví dụ Suggestions

  const handleSearch = (e) => {
    const value = e.target.value;
    setSearchTerm(value);
    setShowResults(value.length > 0);
  };

  return (
    <div className="search-bar-container">
        <div className="input-group search-input-container">
            <div className="input-group-text search-icon-bg border-end-0" style={{ backgroundColor: '#333333',border:1}}>
                <FontAwesomeIcon icon={faSearch} className="text" style={{color:'white'}} />
            </div>
            <input
                type="text"
                className="form-control search-input-field border-start-0"
                placeholder="Search"
                value={searchTerm}
                onChange={handleSearch}
            />
        </div>
        
      
      {showResults && (
        <SearchResults recentSearches={recentSearches} suggestions={suggestions} />
      )}
    </div>
  );
};

export default SearchBar;
