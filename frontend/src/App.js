import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, Routes, Link, Navigate } from 'react-router-dom';
import { 
    AppBar, Toolbar, Typography, Button, Container, TextField, InputAdornment, 
    IconButton, List, ListItem, ListItemText, Badge, Drawer, Box, Grid, Avatar, Divider, Switch 
} from '@mui/material';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import SearchIcon from '@mui/icons-material/Search';
import NotificationsIcon from '@mui/icons-material/Notifications';
import MenuIcon from '@mui/icons-material/Menu';
import Brightness4Icon from '@mui/icons-material/Brightness4';
import Home from './Home';
import Profile from './Profile';
import Photos from './Photos';
import Login from './Login';
import Notifications from './Notifications';
import Groups from './Groups';
import Pages from './Pages';
import FriendRequests from './FriendRequests';
import Explore from './Explore';
import Messaging from './Messaging';
import AIIM from './AIIM';
import Lists from './Lists';
import Hangouts from './Hangouts';
import Settings from './Settings';
import Shorts from './Shorts';
import { getUser, searchUsers } from './api';

function App() {
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [searchQuery, setSearchQuery] = useState('');
    const [searchResults, setSearchResults] = useState([]);
    const [unreadCount, setUnreadCount] = useState(0);
    const [drawerOpen, setDrawerOpen] = useState(false);
    const [user, setUser] = useState(null);
    const [trends, setTrends] = useState(['#Amity', '#SocialMedia', '#EarthTones']);
    const [themeMode, setThemeMode] = useState('light');

    const theme = createTheme({
        palette: {
            mode: themeMode,
            primary: {
                main: '#4a4a4a', // Dark gray
            },
            secondary: {
                main: '#6d6d6d', // Medium gray
            },
            background: {
                default: themeMode === 'light' ? '#e0e0e0' : '#333333', // Light gray / dark gray
                paper: themeMode === 'light' ? '#f5f5f5' : '#424242',   // Off-white / darker gray
            },
            text: {
                primary: themeMode === 'light' ? '#333333' : '#e0e0e0',
                secondary: themeMode === 'light' ? '#555555' : '#b0b0b0',
            },
        },
        typography: {
            fontFamily: 'Roboto, Arial, sans-serif',
            h1: {
                fontSize: '2.5rem',
                fontWeight: 700,
                color: '#4a4a4a',
            },
            h2: {
                fontSize: '1.8rem',
                fontWeight: 600,
                color: '#6d6d6d',
            },
            body1: {
                color: themeMode === 'light' ? '#333333' : '#e0e0e0',
            },
        },
        components: {
            MuiCard: {
                styleOverrides: {
                    root: {
                        borderRadius: 8,
                        boxShadow: '0 4px 8px rgba(0, 0, 0, 0.1)',
                        backgroundColor: themeMode === 'light' ? '#f5f5f5' : '#424242',
                    },
                },
            },
            MuiPaper: {
                styleOverrides: {
                    root: {
                        borderRadius: 8,
                        boxShadow: '0 4px 8px rgba(0, 0, 0, 0.1)',
                        backgroundColor: themeMode === 'light' ? '#f5f5f5' : '#424242',
                    },
                },
            },
            MuiLink: {
                styleOverrides: {
                    root: {
                        color: '#d32f2f', // Red links
                    },
                },
            },
        },
    });

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            setIsAuthenticated(true);
            fetchUser();
        }
    }, []);

    const fetchUser = async () => {
        try {
            const res = await getUser(localStorage.getItem('username'));
            setUser(res.data);
        } catch (err) {
            console.error('Error fetching user:', err);
        }
    };

    const handleLogout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('username');
        setIsAuthenticated(false);
        setUser(null);
    };

    const handleSearch = async () => {
        if (!searchQuery) return;
        try {
            const res = await searchUsers(searchQuery);
            setSearchResults(res.data);
        } catch (err) {
            console.error('Error searching users:', err);
        }
    };

    const toggleDrawer = () => {
        setDrawerOpen(!drawerOpen);
    };

    const toggleTheme = () => {
        setThemeMode(themeMode === 'light' ? 'dark' : 'light');
    };

    const drawerContent = (
        <Box sx={{ width: 250, p: 2 }}>
            {user && (
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                    <Avatar src={user.profile_photo} sx={{ mr: 2 }} />
                    <Typography variant="h6">{user.username}</Typography>
                </Box>
            )}
            <Divider />
            <List>
                {[
                    { text: 'Home', path: '/' },
                    { text: 'Profile', path: '/profile' },
                    { text: 'Photos', path: '/photos' },
                    { text: 'Shorts', path: '/shorts' },
                    { text: 'Notifications', path: '/notifications' },
                    { text: 'Groups', path: '/groups' },
                    { text: 'Pages', path: '/pages' },
                    { text: 'Friend Requests', path: '/friend-requests' },
                    { text: 'Explore', path: '/explore' },
                    { text: 'Messaging', path: '/messaging' },
                    { text: 'AI IM', path: '/ai-im' },
                    { text: 'Lists', path: '/lists' },
                    { text: 'Hangouts', path: '/hangouts' },
                    { text: 'Settings', path: '/settings' },
                ].map(item => (
                    <ListItem button key={item.text} component={Link} to={item.path} onClick={toggleDrawer}>
                        <ListItemText primary={item.text} />
                    </ListItem>
                ))}
                <ListItem button onClick={handleLogout}>
                    <ListItemText primary="Logout" />
                </ListItem>
            </List>
        </Box>
    );

    return (
        <ThemeProvider theme={theme}>
            <Router>
                <AppBar position="static">
                    <Toolbar>
                        <IconButton color="inherit" onClick={toggleDrawer} sx={{ mr: 2 }}>
                            <MenuIcon />
                        </IconButton>
                        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
                            Amity
                        </Typography>
                        {isAuthenticated && (
                            <>
                                <TextField
                                    variant="outlined"
                                    placeholder="Search users..."
                                    value={searchQuery}
                                    onChange={(e) => setSearchQuery(e.target.value)}
                                    onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                                    InputProps={{
                                        endAdornment: (
                                            <InputAdornment position="end">
                                                <IconButton onClick={handleSearch}>
                                                    <SearchIcon />
                                                </IconButton>
                                            </InputAdornment>
                                        ),
                                    }}
                                    sx={{ 
                                        backgroundColor: themeMode === 'light' ? '#f5f5f5' : '#424242', 
                                        borderRadius: 1, 
                                        ml: 2, 
                                        width: { xs: '100%', sm: 'auto' } 
                                    }}
                                    size="small"
                                />
                                <IconButton color="inherit" component={Link} to="/notifications">
                                    <Badge badgeContent={unreadCount} color="secondary">
                                        <NotificationsIcon />
                                    </Badge>
                                </IconButton>
                                <IconButton color="inherit" onClick={toggleTheme}>
                                    <Brightness4Icon />
                                </IconButton>
                            </>
                        )}
                    </Toolbar>
                </AppBar>
                <Drawer anchor="left" open={drawerOpen} onClose={toggleDrawer}>
                    {drawerContent}
                </Drawer>
                <Container sx={{ mt: 4, mb: 4 }}>
                    <Grid container spacing={2}>
                        <Grid item xs={12} md={8}>
                            {isAuthenticated && searchResults.length > 0 && (
                                <List sx={{ mb: 4 }}>
                                    {searchResults.map(result => (
                                        <ListItem key={result.id}>
                                            <ListItemText primary={result.username} secondary={result.location} />
                                        </ListItem>
                                    ))}
                                </List>
                            )}
                            <Routes>
                                <Route path="/login" element={<Login setIsAuthenticated={setIsAuthenticated} />} />
                                <Route 
                                    path="/" 
                                    element={isAuthenticated ? <Home /> : <Navigate to="/login" />} 
                                />
                                <Route 
                                    path="/profile" 
                                    element={isAuthenticated ? <Profile /> : <Navigate to="/login" />} 
                                />
                                <Route 
                                    path="/photos" 
                                    element={isAuthenticated ? <Photos /> : <Navigate to="/login" />} 
                                />
                                <Route 
                                    path="/shorts" 
                                    element={isAuthenticated ? <Shorts /> : <Navigate to="/login" />} 
                                />
                                <Route 
                                    path="/notifications" 
                                    element={isAuthenticated ? <Notifications setUnreadCount={setUnreadCount} /> : <Navigate to="/login" />} 
                                />
                                <Route 
                                    path="/groups" 
                                    element={isAuthenticated ? <Groups /> : <Navigate to="/login" />} 
                                />
                                <Route 
                                    path="/pages" 
                                    element={isAuthenticated ? <Pages /> : <Navigate to="/login" />} 
                                />
                                <Route 
                                    path="/friend-requests" 
                                    element={isAuthenticated ? <FriendRequests /> : <Navigate to="/login" />} 
                                />
                                <Route 
                                    path="/explore" 
                                    element={isAuthenticated ? <Explore /> : <Navigate to="/login" />} 
                                />
                                <Route 
                                    path="/messaging" 
                                    element={isAuthenticated ? <Messaging /> : <Navigate to="/login" />} 
                                />
                                <Route 
                                    path="/ai-im" 
                                    element={isAuthenticated ? <AIIM /> : <Navigate to="/login" />} 
                                />
                                <Route 
                                    path="/lists" 
                                    element={isAuthenticated ? <Lists /> : <Navigate to="/login" />} 
                                />
                                <Route 
                                    path="/hangouts" 
                                    element={isAuthenticated ? <Hangouts /> : <Navigate to="/login" />} 
                                />
                                <Route 
                                    path="/settings" 
                                    element={isAuthenticated ? <Settings /> : <Navigate to="/login" />} 
                                />
                            </Routes>
                        </Grid>
                        <Grid item xs={12} md={4}>
                            {isAuthenticated && (
                                <Box sx={{ p: 2, backgroundColor: 'background.paper', borderRadius: 2, boxShadow: 1 }}>
                                    <Typography variant="h6">Connections</Typography>
                                    <List>
                                        {user?.friends?.map(friend => (
                                            <ListItem key={friend}>
                                                <ListItemText primary={friend} />
                                            </ListItem>
                                        ))}
                                    </List>
                                    <Typography variant="h6" sx={{ mt: 2 }}>Trends</Typography>
                                    <List>
                                        {trends.map(trend => (
                                            <ListItem key={trend}>
                                                <ListItemText primary={trend} />
                                            </ListItem>
                                        ))}
                                    </List>
                                </Box>
                            )}
                        </Grid>
                    </Grid>
                </Container>
            </Router>
        </ThemeProvider>
    );
}

export default App;