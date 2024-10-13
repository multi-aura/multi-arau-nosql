import React from 'react';
import './SearchResults.css';

const SearchResults = ({ recentSearches, suggestions }) => {
  return (
    <div className="search-results-container">
      {/* Recent Section */}
      <div className="recent-searches">
        <h5>Recent</h5>
        <ul>
          {recentSearches.map((item, index) => (
            <li key={index} className="d-flex justify-content-between align-items-center">
              <span className="search-item">{item}</span>
              <button className="btn-remove">X</button>
            </li>
          ))}
        </ul>
        <div className="see-more">See more</div>
      </div>

      <hr />

      {/* Suggestions Section */}
      <div className="suggestions-for-you">
        <h5>Suggestions for you</h5>
        <ul>
          {suggestions.map((item, index) => (
            <li key={index} className="d-flex justify-content-between align-items-center">
              <span className="search-item">{item}</span>
              <button className="btn-remove">X</button>
            </li>
          ))}
        </ul>
        <div className="see-more">See more</div>
      </div>
    </div>
  );
};

export default SearchResults;
