import React from 'react';
import { Container, Grid, TextField, Button, Typography, Box, Link } from '@mui/material';
import './App.css';

function App() {
  return (
    <Container maxWidth="lg" sx={{ mt: 4 }}>
      {/* Header */}
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 4 }}>
        <Typography variant="h4" sx={{ fontWeight: 'bold', color: '#4CAF50' }}>
          GreenHeartPT
        </Typography>
        <Box sx={{ flexGrow: 1 }} />
        <TextField variant="outlined" placeholder="Search" sx={{ width: '300px' }} />
      </Box>

      <Grid container spacing={4}>
        {/* Left Sidebar */}
        <Grid item xs={12} md={4}>
          <Typography variant="h6" gutterBottom>
            Stay connected with everyone.
          </Typography>
          <Typography variant="h6" gutterBottom>
            Promote your post.
          </Typography>
          <Typography variant="h6" gutterBottom>
            Connect with what you like.
          </Typography>
          <Typography variant="h6" gutterBottom>
            Become anonymous.
          </Typography>
          <Typography variant="h6" gutterBottom>
            Create topics & trends.
          </Typography>
          <Typography variant="body2" color="textSecondary" sx={{ mt: 4 }}>
            about policy terms ads promote developer groups language social
          </Typography>
          <Typography variant="body2" color="textSecondary">
            beta: 1.0019 testing in progress!
          </Typography>
          <Typography variant="body2" color="textSecondary">
            An ✕ Aviyon creation
          </Typography>
        </Grid>

        {/* Login Form */}
        <Grid item xs={12} md={8}>
          <Typography variant="h5" gutterBottom>
            Sign up ^_^, it'll always be free!
          </Typography>
          <TextField
            fullWidth
            label="Email"
            variant="outlined"
            sx={{ mb: 2 }}
          />
          <TextField
            fullWidth
            label="Password"
            type="password"
            variant="outlined"
            sx={{ mb: 2 }}
          />
          <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
            <Button variant="contained" color="primary">
              Login
            </Button>
            <Box>
              <Link href="#" sx={{ mr: 2 }}>Create Account</Link>
              <Link href="#">Forgot Something?</Link>
            </Box>
          </Box>

          {/* Stats */}
          <Box sx={{ display: 'flex', alignItems: 'center', mt: 4 }}>
            <Typography variant="body1" sx={{ mr: 2 }}>
              Total Users: 13
            </Typography>
            <Typography variant="body1">
              Total Post: 115
            </Typography>
          </Box>
          {/* Placeholder for post thumbnails */}
          <Box sx={{ display: 'flex', mt: 2 }}>
            <Box sx={{ width: 50, height: 50, bgcolor: 'grey.300', mr: 1 }} />
            <Box sx={{ width: 50, height: 50, bgcolor: 'grey.300', mr: 1 }} />
            <Box sx={{ width: 50, height: 50, bgcolor: 'grey.300' }} />
          </Box>
        </Grid>
      </Grid>
    </Container>
  );
}

export default App;