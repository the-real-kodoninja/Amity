import React, { useState, useEffect } from 'react';
import { Container, Typography, TextField, Button, List, ListItem, ListItemText, Box } from '@mui/material';
import { io } from 'socket.io-client';
import { sendMessage, getMessages } from './api';

const socket = io('http://localhost:8080');

function Messaging() {
    const [messages, setMessages] = useState([]);
    const [newMessage, setNewMessage] = useState('');
    const [recipient, setRecipient] = useState('');

    useEffect(() => {
        socket.on('message', (message) => {
            setMessages(prev => [...prev, message]);
        });

        return () => {
            socket.off('message');
        };
    }, []);

    const fetchMessages = async (username) => {
        try {
            const res = await getMessages(username);
            setMessages(res.data);
            setRecipient(username);
        } catch (err) {
            console.error('Error fetching messages:', err);
        }
    };

    const handleSendMessage = async () => {
        if (!newMessage || !recipient) return;
        try {
            const message = {
                to: recipient,
                content: newMessage,
                is_ai: false
            };
            await sendMessage(message);
            socket.emit('message', { ...message, from: localStorage.getItem('username') });
            setMessages(prev => [...prev, { ...message, from: localStorage.getItem('username'), created_at: new Date().toISOString() }]);
            setNewMessage('');
        } catch (err) {
            console.error('Error sending message:', err);
        }
    };

    return (
        <Container>
            <Typography variant="h4" gutterBottom>Messaging</Typography>
            <TextField
                label="Recipient Username"
                value={recipient}
                onChange={(e) => fetchMessages(e.target.value)}
                fullWidth
                margin="normal"
            />
            <Box sx={{ maxHeight: 400, overflowY: 'auto', mb: 2 }}>
                <List>
                    {messages.map((msg, idx) => (
                        <ListItem key={idx}>
                            <ListItemText
                                primary={`${msg.from}: ${msg.content}`}
                                secondary={msg.created_at}
                                sx={{ textAlign: msg.from === localStorage.getItem('username') ? 'right' : 'left' }}
                            />
                        </ListItem>
                    ))}
                </List>
            </Box>
            <TextField
                label="Message"
                value={newMessage}
                onChange={(e) => setNewMessage(e.target.value)}
                fullWidth
                margin="normal"
            />
            <Button variant="contained" color="primary" onClick={handleSendMessage}>
                Send
            </Button>
        </Container>
    );
}

export default Messaging;