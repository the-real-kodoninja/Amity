import React, { useState, useEffect } from 'react';
import { 
    Container, Typography, TextField, Button, List, ListItem, Box, Grid, Card, CardContent, 
    CardActions, Avatar, IconButton, Menu, MenuItem, CircularProgress, FormControl, InputLabel, Select 
} from '@mui/material';
import { useDropzone } from 'react-dropzone';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import FavoriteIcon from '@mui/icons-material/Favorite';
import ShareIcon from '@mui/icons-material/Share';
import VisibilityOffIcon from '@mui/icons-material/VisibilityOff';
import EmojiEmotionsIcon from '@mui/icons-material/EmojiEmotions';
import GifIcon from '@mui/icons-material/Gif';
import ReactPlayer from 'react-player';
import EmojiPicker from 'emoji-picker-react';
import Giphy from 'giphy-js-sdk-core';
import { format } from 'timeago.js';
import { getFeed, createPost, likePost, reactToPost, sharePost, hidePost, addComment, addReply, uploadFile } from './api';

const giphy = Giphy('your-giphy-api-key'); // Replace with your Giphy API key

function Home() {
    const [posts, setPosts] = useState([]);
    const [newPost, setNewPost] = useState('');
    const [mediaFiles, setMediaFiles] = useState([]);
    const [comment, setComment] = useState({});
    const [reply, setReply] = useState({});
    const [reactionMenu, setReactionMenu] = useState(null);
    const [postMenu, setPostMenu] = useState(null);
    const [selectedPost, setSelectedPost] = useState(null);
    const [emojiPicker, setEmojiPicker] = useState(null);
    const [gifPicker, setGifPicker] = useState(null);
    const [gifs, setGifs] = useState([]);
    const [filter, setFilter] = useState('all');
    const [visiblePosts, setVisiblePosts] = useState(5);
    const [loadingImages, setLoadingImages] = useState({});

    const maxChars = 280;

    useEffect(() => {
        const fetchFeed = async () => {
            try {
                const res = await getFeed(filter);
                setPosts(res.data);
            } catch (err) {
                console.error('Error fetching feed:', err);
            }
        };
        fetchFeed();
    }, [filter]);

    const handleCreatePost = async () => {
        if (newPost.length > maxChars) {
            alert(`Post exceeds ${maxChars} characters`);
            return;
        }

        let media = [];
        for (const file of mediaFiles) {
            const formData = new FormData();
            formData.append('file', file);
            const res = await uploadFile(formData);
            media.push(res.data);
        }

        try {
            await createPost({ 
                content: newPost, 
                username: localStorage.getItem('username'), 
                media, 
                is_short: media.some(m => m.type === 'video' && newPost.length <= 30) 
            });
            setNewPost('');
            setMediaFiles([]);
            const res = await getFeed(filter);
            setPosts(res.data);
        } catch (err) {
            console.error('Error creating post:', err);
        }
    };

    const handleFileUpload = (acceptedFiles) => {
        setMediaFiles([...mediaFiles, ...acceptedFiles]);
    };

    const { getRootProps, getInputProps } = useDropzone({
        accept: 'image/*,video/*,application/*',
        onDrop: handleFileUpload
    });

    const handleLike = async (postId) => {
        try {
            await likePost(postId);
            const res = await getFeed(filter);
            setPosts(res.data);
        } catch (err) {
            console.error('Error liking post:', err);
        }
    };

    const handleReact = async (postId, reaction) => {
        try {
            await reactToPost(postId, reaction);
            const res = await getFeed(filter);
            setPosts(res.data);
        } catch (err) {
            console.error('Error reacting to post:', err);
        }
        setReactionMenu(null);
    };

    const handleShare = async (postId) => {
        try {
            await sharePost(postId);
            const res = await getFeed(filter);
            setPosts(res.data);
        } catch (err) {
            console.error('Error sharing post:', err);
        }
    };

    const handleHide = async (postId) => {
        try {
            await hidePost(postId);
            const res = await getFeed(filter);
            setPosts(res.data.filter(p => !p.hidden_by.includes(localStorage.getItem('username'))));
        } catch (err) {
            console.error('Error hiding post:', err);
        }
        setPostMenu(null);
    };

    const handleComment = async (postId) => {
        try {
            await addComment(postId, { 
                username: localStorage.getItem('username'), 
                content: comment[postId] || '', 
                emoji: comment[`${postId}_emoji`] || '', 
                gif: comment[`${postId}_gif`] || '' 
            });
            setComment({ ...comment, [postId]: '', [`${postId}_emoji`]: '', [`${postId}_gif`]: '' });
            const res = await getFeed(filter);
            setPosts(res.data);
        } catch (err) {
            console.error('Error adding comment:', err);
        }
    };

    const handleReply = async (postId, commentId) => {
        try {
            await addReply(postId, commentId, { 
                username: localStorage.getItem('username'), 
                content: reply[`${postId}_${commentId}`] || '', 
                emoji: reply[`${postId}_${commentId}_emoji`] || '', 
                gif: reply[`${postId}_${commentId}_gif`] || '' 
            });
            setReply({ ...reply, [`${postId}_${commentId}`]: '', [`${postId}_${commentId}_emoji`]: '', [`${postId}_${commentId}_gif`]: '' });
            const res = await getFeed(filter);
            setPosts(res.data);
        } catch (err) {
            console.error('Error adding reply:', err);
        }
    };

    const handleEmojiSelect = (postId, emoji, type = 'comment') => {
        if (type === 'comment') {
            setComment({ ...comment, [`${postId}_emoji`]: emoji.emoji });
        } else {
            setReply({ ...reply, [`${postId}_${type}_emoji`]: emoji.emoji });
        }
        setEmojiPicker(null);
    };

    const handleGifSearch = async (postId, type = 'comment') => {
        try {
            const res = await giphy.search('gifs', { q: 'funny', limit: 10 });
            setGifs(res.data);
            setGifPicker({ postId, type });
        } catch (err) {
            console.error('Error fetching GIFs:', err);
        }
    };

    const handleGifSelect = (postId, gif, type = 'comment') => {
        if (type === 'comment') {
            setComment({ ...comment, [`${postId}_gif`]: gif.images.fixed_height.url });
        } else {
            setReply({ ...reply, [`${postId}_${type}_gif`]: gif.images.fixed_height.url });
        }
        setGifPicker(null);
    };

    const handleImageLoad = (postId, mediaId) => {
        setLoadingImages(prev => ({ ...prev, [`${postId}_${mediaId}`]: false }));
    };

    const renderComments = (comments, postId, level = 0) => (
        comments.map(c => (
            <Box key={c.id} sx={{ ml: level * 2, mt: 1 }}>
                <Box sx={{ display: 'flex', alignItems: 'center' }}>
                    <Avatar sx={{ width: 24, height: 24, mr: 1 }} />
                    <Typography variant="body2">
                        <strong>{c.username}</strong>: {c.content} {c.emoji && <span>{c.emoji}</span>}
                    </Typography>
                </Box>
                {c.gif && <img src={c.gif} alt="GIF" style={{ maxWidth: 200, marginTop: 8 }} />}
                <Typography variant="caption" sx={{ ml: 4 }}>{format(c.timestamp)}</Typography>
                <Button size="small" onClick={() => setReply({ ...reply, [`${postId}_${c.id}_show`]: !reply[`${postId}_${c.id}_show`] })}>
                    Reply
                </Button>
                {reply[`${postId}_${c.id}_show`] && (
                    <Box sx={{ ml: 4, mt: 1 }}>
                        <TextField
                            label="Reply"
                            value={reply[`${postId}_${c.id}`] || ''}
                            onChange={(e) => setReply({ ...reply, [`${postId}_${c.id}`]: e.target.value })}
                            fullWidth
                            size="small"
                        />
                        <IconButton onClick={() => setEmojiPicker({ postId, type: `${c.id}` })}>
                            <EmojiEmotionsIcon />
                        </IconButton>
                        <IconButton onClick={() => handleGifSearch(postId, c.id)}>
                            <GifIcon />
                        </IconButton>
                        <Button onClick={() => handleReply(postId, c.id)}>Reply</Button>
                    </Box>
                )}
                {c.replies && renderComments(c.replies, postId, level + 1)}
            </Box>
        ))
    );

    return (
        <Container>
            <Typography variant="h4" gutterBottom>Home</Typography>
            <Box sx={{ mb: 4 }}>
                <TextField
                    label={`What's on your mind? (${newPost.length}/${maxChars})`}
                    value={newPost}
                    onChange={(e) => setNewPost(e.target.value)}
                    fullWidth
                    margin="normal"
                    error={newPost.length > maxChars}
                    multiline
                    rows={3}
                />
                <Box {...getRootProps()} sx={{ p: 2, border: '1px dashed gray', mb: 2 }}>
                    <input {...getInputProps()} />
                    <Typography>Drag or click to upload photos, videos, or files</Typography>
                </Box>
                {mediaFiles.map((file, idx) => (
                    <Typography key={idx}>{file.name}</Typography>
                ))}
                <Button variant="contained" color="primary" onClick={handleCreatePost}>
                    Post
                </Button>
            </Box>
            <FormControl fullWidth sx={{ mb: 2 }}>
                <InputLabel>Filter Feed</InputLabel>
                <Select value={filter} onChange={(e) => setFilter(e.target.value)}>
                    <MenuItem value="all">All</MenuItem>
                    <MenuItem value="photo">Photos</MenuItem>
                    <MenuItem value="video">Videos</MenuItem>
                    <MenuItem value="file">Files</MenuItem>
                </Select>
            </FormControl>
            <List>
                {posts.slice(0, visiblePosts).map(post => (
                    <ListItem key={post.id}>
                        <Card sx={{ width: '100%' }}>
                            <CardContent>
                                <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                                    <Avatar src={post.profile_photo} sx={{ mr: 2 }} />
                                    <Box sx={{ flexGrow: 1 }}>
                                        <Typography variant="body1">{post.username}</Typography>
                                        <Typography variant="caption">{format(post.timestamp)}</Typography>
                                    </Box>
                                    <IconButton onClick={(e) => { setPostMenu(e.currentTarget); setSelectedPost(post.id); }}>
                                        <MoreVertIcon />
                                    </IconButton>
                                </Box>
                                <Typography variant="body1">{post.content}</Typography>
                                <Grid container spacing={1} sx={{ mt: 2 }}>
                                    {post.media.map((m, idx) => (
                                        <Grid item xs={12} sm={6} key={idx}>
                                            {m.type === 'photo' && (
                                                <>
                                                    {loadingImages[`${post.id}_${idx}`] !== false && (
                                                        <Box sx={{ bgcolor: '#e0e0e0', width: '100%', height: 140 }} />
                                                    )}
                                                    <img
                                                        src={m.url}
                                                        alt="Post media"
                                                        style={{ maxWidth: '100%', display: loadingImages[`${post.id}_${idx}`] === false ? 'block' : 'none' }}
                                                        onLoad={() => handleImageLoad(post.id, idx)}
                                                    />
                                                </>
                                            )}
                                            {m.type === 'video' && (
                                                <ReactPlayer url={m.url} width="100%" height="auto" controls />
                                            )}
                                            {m.type === 'file' && (
                                                <a href={m.url} target="_blank" rel="noopener noreferrer" style={{ color: '#d32f2f' }}>
                                                    Download File ({(m.size / 1024).toFixed(2)} KB)
                                                </a>
                                            )}
                                        </Grid>
                                    ))}
                                </Grid>
                                <Typography variant="caption" display="block" sx={{ mt: 1 }}>
                                    Likes: {post.likes} | Reactions: {Object.entries(post.reactions || {}).map(([k, v]) => `${k}: ${v}`).join(', ')} | Shares: {post.shares}
                                </Typography>
                                {post.comments?.length > 0 && (
                                    <Box sx={{ mt: 2 }}>
                                        {renderComments(post.comments, post.id)}
                                    </Box>
                                )}
                            </CardContent>
                            <CardActions>
                                <IconButton onClick={() => handleLike(post.id)}>
                                    <FavoriteIcon />
                                </IconButton>
                                <IconButton onClick={(e) => { setReactionMenu(e.currentTarget); setSelectedPost(post.id); }}>
                                    <EmojiEmotionsIcon />
                                </IconButton>
                                <IconButton onClick={() => handleShare(post.id)}>
                                    <ShareIcon />
                                </IconButton>
                            </CardActions>
                            <Box sx={{ p: 2 }}>
                                <TextField
                                    label="Add a comment"
                                    value={comment[post.id] || ''}
                                    onChange={(e) => setComment({ ...comment, [post.id]: e.target.value })}
                                    fullWidth
                                    size="small"
                                />
                                <IconButton onClick={() => setEmojiPicker({ postId: post.id, type: 'comment' })}>
                                    <EmojiEmotionsIcon />
                                </IconButton>
                                <IconButton onClick={() => handleGifSearch(post.id)}>
                                    <GifIcon />
                                </IconButton>
                                <Button onClick={() => handleComment(post.id)}>Comment</Button>
                            </Box>
                        </Card>
                    </ListItem>
                ))}
            </List>
            {visiblePosts < posts.length && (
                <Button onClick={() => setVisiblePosts(prev => prev + 5)} sx={{ mt: 2 }}>
                    Load 5 More
                </Button>
            )}
            <Menu
                anchorEl={reactionMenu}
                open={Boolean(reactionMenu)}
                onClose={() => setReactionMenu(null)}
            >
                <MenuItem onClick={() => handleReact(selectedPost, 'heart')}>❤️ Heart</MenuItem>
                <MenuItem onClick={() => handleReact(selectedPost, 'kiss')}>💋 Kiss</MenuItem>
                <MenuItem onClick={() => handleReact(selectedPost, 'laugh')}>😂 Laugh</MenuItem>
            </Menu>
            <Menu
                anchorEl={postMenu}
                open={Boolean(postMenu)}
                onClose={() => setPostMenu(null)}
            >
                <MenuItem onClick={() => handleHide(selectedPost)}>
                    <VisibilityOffIcon sx={{ mr: 1 }} /> Hide Post
                </MenuItem>
            </Menu>
            {emojiPicker && (
                <Box sx={{ position: 'absolute', zIndex: 1000 }}>
                    <EmojiPicker onEmojiClick={(emoji) => handleEmojiSelect(emojiPicker.postId, emoji, emojiPicker.type)} />
                </Box>
            )}
            {gifPicker && (
                <Box sx={{ position: 'absolute', zIndex: 1000, bgcolor: 'white', p: 2, boxShadow: 3 }}>
                    {gifs.map(gif => (
                        <img
                            key={gif.id}
                            src={gif.images.fixed_height.url}
                            alt="GIF"
                            style={{ width: 100, cursor: 'pointer', margin: 4 }}
                            onClick={() => handleGifSelect(gifPicker.postId, gif, gifPicker.type)}
                        />
                    ))}
                </Box>
            )}
        </Container>
    );
}

export default Home;