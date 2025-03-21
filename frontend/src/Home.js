import React, { useState, useEffect } from 'react';
import { Container, Typography, TextField, Button, List, ListItem, ListItemText, Grid, Card, CardContent, CardActions } from '@mui/material';
import { getFeed, createPost, likePost, sharePost, addComment, mintPost } from './api';

function Home() {
    const [posts, setPosts] = useState([]);
    const [newPost, setNewPost] = useState('');
    const [comment, setComment] = useState({});

    useEffect(() => {
        const fetchFeed = async () => {
            try {
                const res = await getFeed();
                setPosts(res.data);
            } catch (err) {
                console.error('Error fetching feed:', err);
            }
        };
        fetchFeed();
    }, []);

    const handleCreatePost = async () => {
        try {
            await createPost({ content: newPost, username: localStorage.getItem('username') });
            setNewPost('');
            const res = await getFeed();
            setPosts(res.data);
        } catch (err) {
            console.error('Error creating post:', err);
        }
    };

    const handleLike = async (postId) => {
        try {
            await likePost(postId);
            const res = await getFeed();
            setPosts(res.data);
        } catch (err) {
            console.error('Error liking post:', err);
        }
    };

    const handleShare = async (postId) => {
        try {
            await sharePost(postId);
            const res = await getFeed();
            setPosts(res.data);
        } catch (err) {
            console.error('Error sharing post:', err);
        }
    };

    const handleComment = async (postId) => {
        try {
            await addComment(postId, { username: localStorage.getItem('username'), content: comment[postId] || '' });
            setComment({ ...comment, [postId]: '' });
            const res = await getFeed();
            setPosts(res.data);
        } catch (err) {
            console.error('Error adding comment:', err);
        }
    };

    const handleMint = async (postId) => {
        try {
            await mintPost(postId);
            alert('Post minted as NFT');
        } catch (err) {
            console.error('Error minting post:', err);
        }
    };

    return (
        <Container>
            <Typography variant="h4" gutterBottom>Home</Typography>
            <TextField
                label="What's on your mind?"
                value={newPost}
                onChange={(e) => setNewPost(e.target.value)}
                fullWidth
                margin="normal"
            />
            <Button variant="contained" color="primary" onClick={handleCreatePost}>
                Post
            </Button>
            <List>
                {posts.map(post => (
                    <ListItem key={post.id}>
                        <Card sx={{ width: '100%' }}>
                            <CardContent>
                                <Typography variant="body1">{post.content}</Typography>
                                <Typography variant="caption">By {post.username}</Typography>
                                <Typography variant="caption" display="block">
                                    Likes: {post.likes} | Shares: {post.shares}
                                </Typography>
                                {post.comments?.map((c, idx) => (
                                    <Typography key={idx} variant="body2">
                                        {c.username}: {c.content}
                                    </Typography>
                                ))}
                            </CardContent>
                            <CardActions>
                                <Button onClick={() => handleLike(post.id)}>Like</Button>
                                <Button onClick={() => handleShare(post.id)}>Share</Button>
                                <Button onClick={() => handleMint(post.id)}>Mint as NFT</Button>
                            </CardActions>
                            <Box sx={{ p: 2 }}>
                                <TextField
                                    label="Add a comment"
                                    value={comment[post.id] || ''}
                                    onChange={(e) => setComment({ ...comment, [post.id]: e.target.value })}
                                    fullWidth
                                />
                                <Button onClick={() => handleComment(post.id)}>Comment</Button>
                            </Box>
                        </Card>
                    </ListItem>
                ))}
            </List>
        </Container>
    );
}

export default Home;