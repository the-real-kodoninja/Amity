import React, { useState, useEffect } from 'react';
import { Container, Typography, FormControlLabel, Switch, FormControl, InputLabel, Select, MenuItem, Button } from '@mui/material';
import { getUser, updateUser } from './api';

function Settings() {
    const [settings, setSettings] = useState({
        nsfw_enabled: false,
        privacy: { profile_visibility: 'public', messaging: 'everyone' }
    });

    useEffect(() => {
        const fetchSettings = async () => {
            try {
                const res = await getUser(localStorage.getItem('username'));
                setSettings(res.data.settings);
            } catch (err) {
                console.error('Error fetching settings:', err);
            }
        };
        fetchSettings();
    }, []);

    const handleSave = async () => {
        try {
            await updateUser(localStorage.getItem('username'), { settings });
            alert('Settings updated');
        } catch (err) {
            console.error('Error updating settings:', err);
        }
    };

    return (
        <Container>
            <Typography variant="h4" gutterBottom>Settings</Typography>
            <FormControlLabel
                control={
                    <Switch
                        checked={settings.nsfw_enabled}
                        onChange={(e) => setSettings({ ...settings, nsfw_enabled: e.target.checked })}
                    />
                }
                label="Show NSFW Content"
            />
            <FormControl fullWidth margin="normal">
                <InputLabel>Profile Visibility</InputLabel>
                <Select
                    value={settings.privacy.profile_visibility}
                    onChange={(e) => setSettings({ ...settings, privacy: { ...settings.privacy, profile_visibility: e.target.value } })}
                >
                    <MenuItem value="public">Public</MenuItem>
                    <MenuItem value="friends">Friends Only</MenuItem>
                    <MenuItem value="private">Private</MenuItem>
                </Select>
            </FormControl>
            <FormControl fullWidth margin="normal">
                <InputLabel>Messaging</InputLabel>
                <Select
                    value={settings.privacy.messaging}
                    onChange={(e) => setSettings({ ...settings, privacy: { ...settings.privacy, messaging: e.target.value } })}
                >
                    <MenuItem value="everyone">Everyone</MenuItem>
                    <MenuItem value="friends">Friends Only</MenuItem>
                    <MenuItem value="none">No One</MenuItem>
                </Select>
            </FormControl>
            <Button variant="contained" color="primary" onClick={handleSave}>
                Save Settings
            </Button>
        </Container>
    );
}

export default Settings;