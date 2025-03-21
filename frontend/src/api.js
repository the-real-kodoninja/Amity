import axios from 'axios';

const api = axios.create({
    baseURL: '/api',
    headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
    }
});

export const login = (username, password) => api.post('/login', { username, password });
export const register = (userData) => api.post('/register', userData);
export const getUser = (username) => api.get(`/users/${username}`);
export const updateUser = (username, updates) => api.put(`/users/${username}/update`, updates);
export const getFeed = (filter) => api.get(`/feed${filter ? `?filter=${filter}` : ''}`);
export const getExplore = () => api.get('/explore');
export const getUserPhotos = (username) => api.get(`/users/${username}/photos`);
export const createPost = (post) => api.post('/posts', post);
export const likePost = (postId) => api.post(`/posts/${postId}/like`);
export const reactToPost = (postId, reaction) => api.post(`/posts/${postId}/react`, { type: reaction });
export const sharePost = (postId) => api.post(`/posts/${postId}/share`);
export const hidePost = (postId) => api.post(`/posts/${postId}/hide`);
export const addComment = (postId, comment) => api.post(`/posts/${postId}/comment`, comment);
export const addReply = (postId, commentId, reply) => api.post(`/posts/${postId}/comment/${commentId}/reply`, reply);
export const getShorts = () => api.get('/shorts');
export const uploadFile = (formData) => api.post('/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
});

export default api;