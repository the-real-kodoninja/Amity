import React, { useState, useEffect } from 'react';
import { Container, Typography, TextField, Button, List, ListItem, ListItemText } from '@mui/material';
import { createHangout, getHangouts, joinHangout, leaveHangout } from './api';

function Hangouts() {
    const [hangouts, setHangouts] = useState([]);
    const [newHangout, setNewHangout] = useState({ name: '', description: '', date: '', location: '' });

    useEffect(() => {
        const fetchHangouts = async () => {
            try {
                const res = await getHangouts();
                setHangouts(res.data);
            } catch (err) {
                console.error('Error fetching hangouts:', err);
            }
        };
        fetchHangouts();
    }, []);

    const handleCreateHangout = async () => {
        try {
            await createHangout(newHangout);
            setNewHangout({ name: '', description: '', date: '', location: '' });
            const res = await getHangouts();
            setHangouts(res.data);
        } catch (err) {
            console.error('Error creating hangout:', err);
        }
    };

    const handleJoinHangout = async (hangoutId) => {
        try {
            await joinHangout(hangoutId);
            const res = await getHangouts();
            setHangouts(res.data);
        } catch (err) {
            console.error('Error joining hangout:', err);
        }
    };

    const handleLeaveHangout = async (hangoutId) => {
        try {
            await leaveHangout(hangoutId);
            const res = await getHangouts();
            setHangouts(res.data);
        } catch (err) {
            console.error('Error leaving hangout:', err);
        }
    };

    return (
        <Container>
            <Typography variant="h4" gutterBottom>Hangouts</Typography>
            <TextField
                label="Hangout Name"
                value={newHangout.name}
                onChange={(e) => setNewHangout({ ...newHangout, name: e.target.value })}
                fullWidth
                margin="normal"
            />
            <TextField
                label="Description"
                value={newHangout.description}
                onChange={(e) => setNewHangout({ ...newHangout, description: e.target.value })}
                fullWidth
                margin="normal"
            />
            <TextField
                label="Date"
                type="datetime-local"
                value={newHangout.date}
                onChange={(e) => setNewHangout({ ...newHangout, date: e.target.value })}
                fullWidth
                margin="normal"
                InputLabelProps={{ shrink: true }}
            />
            <TextField
                label="Location"
                value={newHangout.location}
                onChange={(e) => setNewHangout({ ...newHangout, location: e.target.value })}
                fullWidth
                margin="normal"
            />
            <Button variant="contained" color="primary" onClick={handleCreateHangout}>
                Create Hangout
            </Button>
            <List>
                {hangouts.map(hangout => (
                    <ListItem key={hangout.id}>
                        <ListItemText
                            primary={hangout.name}
                            secondary={`Date: ${hangout.date}, Location: ${hangout.location}`}
                        />
                        {hangout.participants.includes(localStorage.getItem('username')) ? (
                            <Button onClick={() => handleLeaveHangout(hangout.id)}>Leave</Button>
                        ) : (
                            <Button onClick={() => handleJoinHangout(hangout.id)}>Join</Button>
                        )}
                    </ListItem>
                ))}
            </List>
        </Container>
    );
}

export default Hangouts;