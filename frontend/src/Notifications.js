import React, { useState, useEffect } from 'react';
import { Container, Typography, List, ListItem, ListItemText, Button } from '@mui/material';
import { getNotifications, markNotificationRead } from './api';

function Notifications({ setUnreadCount }) {
    const [notifications, setNotifications] = useState([]);

    useEffect(() => {
        const fetchNotifications = async () => {
            try {
                const res = await getNotifications();
                setNotifications(res.data);
                const unread = res.data.filter(notif => !notif.is_read).length;
                setUnreadCount(unread);
            } catch (err) {
                console.error('Error fetching notifications:', err);
            }
        };
        fetchNotifications();
    }, [setUnreadCount]);

    const handleMarkRead = async (notifId) => {
        try {
            await markNotificationRead(notifId);
            setNotifications(notifications.map(notif =>
                notif.id === notifId ? { ...notif, is_read: true } : notif
            ));
            setUnreadCount(notifications.filter(notif => !notif.is_read && notif.id !== notifId).length);
        } catch (err) {
            console.error('Error marking notification as read:', err);
        }
    };

    return (
        <Container>
            <Typography variant="h4" gutterBottom>Notifications</Typography>
            <List>
                {notifications.map(notif => (
                    <ListItem key={notif.id}>
                        <ListItemText
                            primary={notif.message}
                            secondary={`From: ${notif.from} at ${notif.created_at}`}
                            sx={{ color: notif.is_read ? 'text.secondary' : 'text.primary' }}
                        />
                        {!notif.is_read && (
                            <Button onClick={() => handleMarkRead(notif.id)}>Mark as Read</Button>
                        )}
                    </ListItem>
                ))}
            </List>
        </Container>
    );
}

export default Notifications;