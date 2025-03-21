import React, { useState, useEffect } from 'react';
import { Container, Typography, TextField, Button, List, ListItem, ListItemText, Box } from '@mui/material';
import { io } from 'socket.io-client';
import { sendMessage } from './api';

const socket = io('http://localhost:8080');

function AIIM() {
    const [messages, setMessages] = useState([]);
    const [newMessage, setNewMessage] = useState('');

    useEffect(() => {
        socket.on('message', (message) => {
            if (message.is_ai) {
                setMessages(prev => [...prev, message]);
            }
        });

        return () => {
            socket.off('message');
        };
    }, []);

    const handleSendMessage = async () => {
        if (!newMessage) return;
        try {
            const message = {
                to: 'AI Assistant',
                content: newMessage,
                is_ai: true
            };
            await sendMessage(message);
            setMessages(prev => [...prev, { ...message, from: localStorage.getItem('username'), created_at: new Date().toISOString() }]);
            socket.emit('message', { ...message, from: localStorage.getItem('username') });
            setNewMessage('');
        } catch (err) {
            console.error('Error sending AI message:', err);
        }
    };

    return (
        <Container>
            <Typography variant="h4" gutterBottom>AI Instant Messaging</Typography>
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
                label="Message to AI"
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

export default AIIM;