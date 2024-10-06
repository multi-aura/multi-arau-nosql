import React from 'react';
import Post from '../Post/Post';
import avatar1 from '../../assets/9ea1c9841cadbef3e7bc.jpg';
import avatar2 from '../../assets/Dinkaru.jpg';
import './Feed.css';
const posts = [
  {
    id: 1,
    username: 'Nguyen Huy Hoang',
    avatar: avatar1,
    time: 'Now',
    content: 'Dưới hình bật chú',
    image: avatar1, 
    comments: [
      { username: 'Kim Định', content: 'Tôi đoán là bí và củ dền' },
      { username: 'Nguyen Huy Hoang', content: 'Sai nha' },
    ],
  },
  {
    id: 2,
    username: 'Nguyen Huy Hoang',
    avatar: avatar1,
    time: '2 mins ago',
    content: 'Gấu leo cây',
    image: avatar2,
    likes: 54,
    comments: [
      { username: 'Kim Định', content: 'Hình ảnh dễ thương!' },
      { username: 'Nguyen Huy Hoang', content: 'UwUwwwww' },
    ],
  },
];

function Feed() {
  return (
    <div className="feed">
      {posts.map(post => (
        <Post key={post.id} post={post} />
      ))}
    </div>
  );
}

export default Feed;
