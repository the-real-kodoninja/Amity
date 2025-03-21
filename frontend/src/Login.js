import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { Container, Typography, TextField, Button, Paper } from '@mui/material';

function Login({ setIsAuthenticated }) {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const handleLogin = async (e) => {
        e.preventDefault();
        try {
            const res = await axios.post('/login', { username, password });
            localStorage.setItem('token', res.data.token);
            setIsAuthenticated(true);
            navigate('/');
        } catch (err) {
            setError('Invalid username or password');
        }
    };

    return (
        <Container sx={{ mt: 4, mb: 4 }}>
            <Paper elevation={3} sx={{ p: 3, maxWidth: 400, mx: 'auto' }}>
                <Typography variant="h2" gutterBottom>Login</Typography>
                {error && <Typography color="error">{error}</Typography>}
                <form onSubmit={handleLogin}>
                    <TextField
                        fullWidth
                        label="Username"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        variant="outlined"
                        sx={{ mb: 2 }}
                    />
                    <TextField
                        fullWidth
                        label="Password"
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        variant="outlined"
                        sx={{ mb: 2 }}
                    />
                    <Button type="submit" variant="contained" color="primary" fullWidth>
                        Login
                    </Button>
                </form>
            </Paper>
        </Container>
    );
}

export default Login;