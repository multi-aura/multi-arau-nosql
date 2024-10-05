import React from 'react';
import Comment from '../Comment/Comment';
import './Post.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faThumbsUp, faCommentDots, faShare, faHeart, faBookmark } from '@fortawesome/free-solid-svg-icons'; // Import icon bookmark

function Post({ post }) {
  return (
    <div className="post p-3 mb-4 rounded shadow-sm text-white">
      <div className="d-flex align-items-center mb-3">
        <img src={post.avatar} alt="Avatar" className="avatar rounded-circle" />
        <div className="ml-3">
          <h5>{post.username}</h5>
          <p className="text-muted">{post.time}</p> 
        </div>
      </div>
      <p>{post.content}</p>
      <img src={post.image} alt="Post" className="img-post img-fluid rounded mb-4" />
      <div className="d-flex justify-content-between align-items-center">
        <div className="d-flex">
          <button className="btn btn-link text-white mr-3">
            <FontAwesomeIcon icon={faHeart} />
          </button>
          <button className="btn btn-link text-white mr-3">
            <FontAwesomeIcon icon={faCommentDots} /> 
          </button>
          <button className="btn btn-link text-white mr-3">
            <FontAwesomeIcon icon={faShare} /> 
          </button>
        </div>
        <button className="btn btn-link text-white">
          <FontAwesomeIcon icon={faBookmark} /> 
        </button>
      </div>
      <div className="comments mt-3">
        {post.comments.map((comment, index) => (
          <Comment key={index} comment={comment} />
        ))}
      </div>
      <div className="d-flex mt-3">
        <input type="text" className="form-control" placeholder="Add a comment..." />
        <button className="btn btn-primary ml-2">Post</button>
      </div>
    </div>
  );
}

export default Post;
