import React, { useState, useEffect } from 'react';
import { Container, Typography, TextField, Button, List, ListItem, ListItemText } from '@mui/material';
import { createList, getLists, addToList } from './api';

function Lists() {
    const [lists, setLists] = useState([]);
    const [newListName, setNewListName] = useState('');
    const [itemToAdd, setItemToAdd] = useState({ type: '', item_id: '' });

    useEffect(() => {
        const fetchLists = async () => {
            try {
                const res = await getLists();
                setLists(res.data);
            } catch (err) {
                console.error('Error fetching lists:', err);
            }
        };
        fetchLists();
    }, []);

    const handleCreateList = async () => {
        try {
            await createList({ name: newListName });
            setNewListName('');
            const res = await getLists();
            setLists(res.data);
        } catch (err) {
            console.error('Error creating list:', err);
        }
    };

    const handleAddToList = async (listId) => {
        try {
            await addToList(listId, itemToAdd);
            setItemToAdd({ type: '', item_id: '' });
            const res = await getLists();
            setLists(res.data);
        } catch (err) {
            console.error('Error adding to list:', err);
        }
    };

    return (
        <Container>
            <Typography variant="h4" gutterBottom>Lists</Typography>
            <TextField
                label="New List Name"
                value={newListName}
                onChange={(e) => setNewListName(e.target.value)}
                fullWidth
                margin="normal"
            />
            <Button variant="contained" color="primary" onClick={handleCreateList}>
                Create List
            </Button>
            <Typography variant="h6" sx={{ mt: 4 }}>Add to List</Typography>
            <TextField
                label="Item Type (e.g., post, photo)"
                value={itemToAdd.type}
                onChange={(e) => setItemToAdd({ ...itemToAdd, type: e.target.value })}
                fullWidth
                margin="normal"
            />
            <TextField
                label="Item ID"
                value={itemToAdd.item_id}
                onChange={(e) => setItemToAdd({ ...itemToAdd, item_id: e.target.value })}
                fullWidth
                margin="normal"
            />
            <List>
                {lists.map(list => (
                    <ListItem key={list.id}>
                        <ListItemText primary={list.name} secondary={`Items: ${list.items.length}`} />
                        <Button onClick={() => handleAddToList(list.id)}>Add to This List</Button>
                    </ListItem>
                ))}
            </List>
        </Container>
    );
}

export default Lists;