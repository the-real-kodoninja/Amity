import React, { useState, useEffect } from 'react';
import { Container, Typography, Button, TextField, List, ListItem, ListItemText } from '@mui/material';
import { getPages, createPage, followPage, unfollowPage } from './api';

function Pages() {
    const [pages, setPages] = useState([]);
    const [newPage, setNewPage] = useState({ name: '', description: '' });

    useEffect(() => {
        const fetchPages = async () => {
            try {
                const res = await getPages();
                setPages(res.data);
            } catch (err) {
                console.error('Error fetching pages:', err);
            }
        };
        fetchPages();
    }, []);

    const handleCreatePage = async () => {
        try {
            await createPage(newPage);
            setNewPage({ name: '', description: '' });
            const res = await getPages();
            setPages(res.data);
        } catch (err) {
            console.error('Error creating page:', err);
        }
    };

    const handleFollowPage = async (pageId) => {
        try {
            await followPage(pageId);
            const res = await getPages();
            setPages(res.data);
        } catch (err) {
            console.error('Error following page:', err);
        }
    };

    const handleUnfollowPage = async (pageId) => {
        try {
            await unfollowPage(pageId);
            const res = await getPages();
            setPages(res.data);
        } catch (err) {
            console.error('Error unfollowing page:', err);
        }
    };

    return (
        <Container>
            <Typography variant="h4" gutterBottom>Pages</Typography>
            <TextField
                label="Page Name"
                value={newPage.name}
                onChange={(e) => setNewPage({ ...newPage, name: e.target.value })}
                fullWidth
                margin="normal"
            />
            <TextField
                label="Description"
                value={newPage.description}
                onChange={(e) => setNewPage({ ...newPage, description: e.target.value })}
                fullWidth
                margin="normal"
            />
            <Button variant="contained" color="primary" onClick={handleCreatePage}>
                Create Page
            </Button>
            <List>
                {pages.map(page => (
                    <ListItem key={page.id}>
                        <ListItemText primary={page.name} secondary={page.description} />
                        {page.followers.includes(localStorage.getItem('username')) ? (
                            <Button onClick={() => handleUnfollowPage(page.id)}>Unfollow</Button>
                        ) : (
                            <Button onClick={() => handleFollowPage(page.id)}>Follow</Button>
                        )}
                    </ListItem>
                ))}
            </List>
        </Container>
    );
}

export default Pages;