import React, { useState, useEffect } from 'react';
import { Container, Typography, Grid, Card, CardMedia, CardContent, Button } from '@mui/material';
import { getUserPhotos, mintPhoto } from './api';

function Photos() {
    const [photos, setPhotos] = useState([]);
    const username = localStorage.getItem('username');

    useEffect(() => {
        const fetchPhotos = async () => {
            try {
                const res = await getUserPhotos(username);
                setPhotos(res.data);
            } catch (err) {
                console.error('Error fetching photos:', err);
            }
        };
        fetchPhotos();
    }, [username]);

    const handleMint = async (photoId) => {
        try {
            await mintPhoto(photoId);
            alert('Photo minted as NFT');
        } catch (err) {
            console.error('Error minting photo:', err);
        }
    };

    return (
        <Container>
            <Typography variant="h4" gutterBottom>Photos</Typography>
            <Grid container spacing={2}>
                {photos.map(photo => (
                    <Grid item xs={12} sm={6} md={4} key={photo.id}>
                        <Card>
                            <CardMedia component="img" height="140" image={photo.url} alt="User photo" />
                            <CardContent>
                                <Typography variant="caption">{photo.caption}</Typography>
                            </CardContent>
                            <Button onClick={() => handleMint(photo.id)}>Mint as NFT</Button>
                        </Card>
                    </Grid>
                ))}
            </Grid>
        </Container>
    );
}

export default Photos;