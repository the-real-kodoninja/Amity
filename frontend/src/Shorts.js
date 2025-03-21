import React, { useState, useEffect } from 'react';
import { Container, Typography, Box } from '@mui/material';
import ReactPlayer from 'react-player';
import { getShorts } from './api';

function Shorts() {
    const [shorts, setShorts] = useState([]);

    useEffect(() => {
        const fetchShorts = async () => {
            try {
                const res = await getShorts();
                setShorts(res.data);
            } catch (err) {
                console.error('Error fetching shorts:', err);
            }
        };
        fetchShorts();
    }, []);

    return (
        <Container>
            <Typography variant="h4" gutterBottom>Shorts</Typography>
            {shorts.map(short => (
                <Box key={short.id} sx={{ mb: 4, maxWidth: 400, mx: 'auto' }}>
                    <ReactPlayer
                        url={short.media.find(m => m.type === 'video')?.url}
                        width="100%"
                        height="600px"
                        controls
                        playing
                        loop
                    />
                    <Typography variant="body1">{short.content}</Typography>
                    <Typography variant="caption">By {short.username} • {short.timestamp}</Typography>
                </Box>
            ))}
        </Container>
    );
}

export default Shorts;