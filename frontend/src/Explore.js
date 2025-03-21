import React, { useState, useEffect } from 'react';
import { Container, Typography, List, ListItem, ListItemText, Grid, Card, CardMedia, CardContent } from '@mui/material';
import { getExplore } from './api';

function Explore() {
    const [exploreData, setExploreData] = useState({ posts: [], groups: [], pages: [], users: [] });

    useEffect(() => {
        const fetchExplore = async () => {
            try {
                const res = await getExplore();
                setExploreData(res.data);
            } catch (err) {
                console.error('Error fetching explore data:', err);
            }
        };
        fetchExplore();
    }, []);

    return (
        <Container>
            <Typography variant="h4" gutterBottom>Explore</Typography>
            <Typography variant="h6">Trending Posts</Typography>
            <Grid container spacing={2}>
                {exploreData.posts.map(post => (
                    <Grid item xs={12} sm={6} md={4} key={post.id}>
                        <Card>
                            <CardContent>
                                <Typography variant="body1">{post.content}</Typography>
                                <Typography variant="caption">By {post.username}</Typography>
                            </CardContent>
                        </Card>
                    </Grid>
                ))}
            </Grid>
            <Typography variant="h6" sx={{ mt: 4 }}>Groups</Typography>
            <List>
                {exploreData.groups.map(group => (
                    <ListItem key={group.id}>
                        <ListItemText primary={group.name} secondary={group.description} />
                    </ListItem>
                ))}
            </List>
            <Typography variant="h6" sx={{ mt: 4 }}>Pages</Typography>
            <List>
                {exploreData.pages.map(page => (
                    <ListItem key={page.id}>
                        <ListItemText primary={page.name} secondary={page.description} />
                    </ListItem>
                ))}
            </List>
            <Typography variant="h6" sx={{ mt: 4 }}>Popular Users</Typography>
            <List>
                {exploreData.users.map(user => (
                    <ListItem key={user.id}>
                        <ListItemText primary={user.username} secondary={`Followers: ${user.followers}`} />
                    </ListItem>
                ))}
            </List>
        </Container>
    );
}

export default Explore;