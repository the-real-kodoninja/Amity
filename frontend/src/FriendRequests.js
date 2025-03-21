import React, { useState, useEffect } from 'react';
import { Container, Typography, List, ListItem, ListItemText, Button } from '@mui/material';
import { getFriendRequests, acceptFriendRequest, rejectFriendRequest } from './api';

function FriendRequests() {
    const [requests, setRequests] = useState([]);

    useEffect(() => {
        const fetchRequests = async () => {
            try {
                const res = await getFriendRequests();
                setRequests(res.data);
            } catch (err) {
                console.error('Error fetching friend requests:', err);
            }
        };
        fetchRequests();
    }, []);

    const handleAccept = async (requestId) => {
        try {
            await acceptFriendRequest(requestId);
            setRequests(requests.filter(req => req.id !== requestId));
        } catch (err) {
            console.error('Error accepting friend request:', err);
        }
    };

    const handleReject = async (requestId) => {
        try {
            await rejectFriendRequest(requestId);
            setRequests(requests.filter(req => req.id !== requestId));
        } catch (err) {
            console.error('Error rejecting friend request:', err);
        }
    };

    return (
        <Container>
            <Typography variant="h4" gutterBottom>Friend Requests</Typography>
            <List>
                {requests.map(req => (
                    <ListItem key={req.id}>
                        <ListItemText primary={`From: ${req.from}`} />
                        <Button onClick={() => handleAccept(req.id)} color="primary">Accept</Button>
                        <Button onClick={() => handleReject(req.id)} color="secondary">Reject</Button>
                    </ListItem>
                ))}
            </List>
        </Container>
    );
}

export default FriendRequests;