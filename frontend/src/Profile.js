import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Container, Typography, Button, Box, Tabs, Tab, Avatar, Grid, Card, CardMedia, CardContent, List, ListItem, ListItemText } from '@mui/material';
import { useDropzone } from 'react-dropzone';
import { getUser, updateUser, followUser, unfollowUser, blockUser, unblockUser } from './api';

function Profile() {
    const navigate = useNavigate();
    const [user, setUser] = useState(null);
    const [tabValue, setTabValue] = useState(0);
    const username = localStorage.getItem('username');
    const isOwnProfile = true; // In a real app, compare with URL param

    useEffect(() => {
        const fetchUser = async () => {
            try {
                const res = await getUser(username);
                setUser(res.data);
            } catch (err) {
                console.error('Error fetching user:', err);
            }
        };
        fetchUser();
    }, [username]);

    const handleTabChange = (event, newValue) => {
        setTabValue(newValue);
    };

    const handleFollow = async () => {
        try {
            await followUser(username);
            const res = await getUser(username);
            setUser(res.data);
        } catch (err) {
            console.error('Error following user:', err);
        }
    };

    const handleUnfollow = async () => {
        try {
            await unfollowUser(username);
            const res = await getUser(username);
            setUser(res.data);
        } catch (err) {
            console.error('Error unfollowing user:', err);
        }
    };

    const handleBlock = async () => {
        try {
            await blockUser(username);
            const res = await getUser(username);
            setUser(res.data);
        } catch (err) {
            console.error('Error blocking user:', err);
        }
    };

    const handleUnblock = async () => {
        try {
            await unblockUser(username);
            const res = await getUser(username);
            setUser(res.data);
        } catch (err) {
            console.error('Error unblocking user:', err);
        }
    };

    const handleProfilePhotoUpload = async (acceptedFiles) => {
        const file = acceptedFiles[0];
        const formData = new FormData();
        formData.append('file', file);
        // In a real app, upload to a server and get URL
        const photoUrl = URL.createObjectURL(file);
        try {
            await updateUser(username, { profile_photo: photoUrl });
            const res = await getUser(username);
            setUser(res.data);
        } catch (err) {
            console.error('Error uploading profile photo:', err);
        }
    };

    const handleBannerUpload = async (acceptedFiles) => {
        const file = acceptedFiles[0];
        const formData = new FormData();
        formData.append('file', file);
        const bannerUrl = URL.createObjectURL(file);
        try {
            await updateUser(username, { banner: bannerUrl });
            const res = await getUser(username);
            setUser(res.data);
        } catch (err) {
            console.error('Error uploading banner:', err);
        }
    };

    const profileDropzone = useDropzone({
        accept: 'image/*',
        onDrop: handleProfilePhotoUpload
    });

    const bannerDropzone = useDropzone({
        accept: 'image/*',
        onDrop: handleBannerUpload
    });

    if (!user) return <Typography>Loading...</Typography>;

    return (
        <Container>
            <Box sx={{ position: 'relative', height: 200 }}>
                <Box
                    sx={{
                        height: '100%',
                        backgroundImage: `url(${user.banner})`,
                        backgroundSize: 'cover',
                        backgroundPosition: 'center'
                    }}
                    {...bannerDropzone.getRootProps()}
                >
                    <input {...bannerDropzone.getInputProps()} />
                    {isOwnProfile && (
                        <Typography sx={{ p: 2, color: 'white', backgroundColor: 'rgba(0,0,0,0.5)' }}>
                            Drag or click to upload banner
                        </Typography>
                    )}
                </Box>
                <Avatar
                    src={user.profile_photo}
                    sx={{
                        width: 120,
                        height: 120,
                        position: 'absolute',
                        bottom: -60,
                        left: 20,
                        border: '4px solid white'
                    }}
                    {...profileDropzone.getRootProps()}
                >
                    <input {...profileDropzone.getInputProps()} />
                    {isOwnProfile && !user.profile_photo && 'Upload'}
                </Avatar>
            </Box>
            <Box sx={{ mt: 8, mb: 4 }}>
                <Typography variant="h4">{user.username}</Typography>
                <Typography variant="body1">{user.location}</Typography>
                {!isOwnProfile && (
                    <Box sx={{ mt: 2 }}>
                        {user.following.includes(username) ? (
                            <Button onClick={handleUnfollow}>Unfollow</Button>
                        ) : (
                            <Button onClick={handleFollow}>Follow</Button>
                        )}
                        {user.blocked_users.includes(username) ? (
                            <Button onClick={handleUnblock} color="secondary">Unblock</Button>
                        ) : (
                            <Button onClick={handleBlock} color="secondary">Block</Button>
                        )}
                    </Box>
                )}
            </Box>
            <Tabs value={tabValue} onChange={handleTabChange} centered>
                <Tab label={`Photos (${user.photos?.length || 0})`} />
                <Tab label={`Connections (${user.friends?.length || 0})`} />
                <Tab label={`Followers (${user.followers || 0})`} />
            </Tabs>
            {tabValue === 0 && (
                <Grid container spacing={2} sx={{ mt: 2 }}>
                    {user.photos?.map(photo => (
                        <Grid item xs={12} sm={6} md={4} key={photo.id}>
                            <Card>
                                <CardMedia component="img" height="140" image={photo.url} alt="User photo" />
                                <CardContent>
                                    <Typography variant="caption">{photo.caption}</Typography>
                                </CardContent>
                            </Card>
                        </Grid>
                    ))}
                </Grid>
            )}
            {tabValue === 1 && (
                <List>
                    {user.friends?.map(friend => (
                        <ListItem key={friend}>
                            <ListItemText primary={friend} />
                        </ListItem>
                    ))}
                </List>
            )}
            {tabValue === 2 && (
                <Typography variant="body1" sx={{ mt: 2 }}>
                    Followers: {user.followers}
                </Typography>
            )}
        </Container>
    );
}

export default Profile;