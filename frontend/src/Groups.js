import React, { useState, useEffect } from 'react';
import { Container, Typography, Button, TextField, List, ListItem, ListItemText } from '@mui/material';
import { getGroups, createGroup, joinGroup, leaveGroup } from './api';

function Groups() {
    const [groups, setGroups] = useState([]);
    const [newGroup, setNewGroup] = useState({ name: '', description: '' });

    useEffect(() => {
        const fetchGroups = async () => {
            try {
                const res = await getGroups();
                setGroups(res.data);
            } catch (err) {
                console.error('Error fetching groups:', err);
            }
        };
        fetchGroups();
    }, []);

    const handleCreateGroup = async () => {
        try {
            await createGroup(newGroup);
            setNewGroup({ name: '', description: '' });
            const res = await getGroups();
            setGroups(res.data);
        } catch (err) {
            console.error('Error creating group:', err);
        }
    };

    const handleJoinGroup = async (groupId) => {
        try {
            await joinGroup(groupId);
            const res = await getGroups();
            setGroups(res.data);
        } catch (err) {
            console.error('Error joining group:', err);
        }
    };

    const handleLeaveGroup = async (groupId) => {
        try {
            await leaveGroup(groupId);
            const res = await getGroups();
            setGroups(res.data);
        } catch (err) {
            console.error('Error leaving group:', err);
        }
    };

    return (
        <Container>
            <Typography variant="h4" gutterBottom>Groups</Typography>
            <TextField
                label="Group Name"
                value={newGroup.name}
                onChange={(e) => setNewGroup({ ...newGroup, name: e.target.value })}
                fullWidth
                margin="normal"
            />
            <TextField
                label="Description"
                value={newGroup.description}
                onChange={(e) => setNewGroup({ ...newGroup, description: e.target.value })}
                fullWidth
                margin="normal"
            />
            <Button variant="contained" color="primary" onClick={handleCreateGroup}>
                Create Group
            </Button>
            <List>
                {groups.map(group => (
                    <ListItem key={group.id}>
                        <ListItemText primary={group.name} secondary={group.description} />
                        {group.members.includes(localStorage.getItem('username')) ? (
                            <Button onClick={() => handleLeaveGroup(group.id)}>Leave</Button>
                        ) : (
                            <Button onClick={() => handleJoinGroup(group.id)}>Join</Button>
                        )}
                    </ListItem>
                ))}
            </List>
        </Container>
    );
}

export default Groups;