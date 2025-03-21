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
export const searchUsers = (query) => api.get(`/search/users`, { params: { q: query } });
export const followUser = (username) => api.post(`/users/${username}/follow`);
export const unfollowUser = (username) => api.post(`/users/${username}/unfollow`);
export const blockUser = (username) => api.post(`/users/${username}/block`);
export const unblockUser = (username) => api.post(`/users/${username}/unblock`);
export const mintPhoto = (photoId) => api.post(`/photos/${photoId}/mint`);
export const getNotifications = () => api.get(`/notifications`);
export const markNotificationRead = (notificationId) => api.post(`/notifications/${notificationId}/read`);
export const getGroups = () => api.get(`/groups`);
export const createGroup = (groupData) => api.post(`/groups`, groupData);
export const joinGroup = (groupId) => api.post(`/groups/${groupId}/join`);
export const leaveGroup = (groupId) => api.post(`/groups/${groupId}/leave`);
export const getFriendRequests = () => api.get(`/friend-requests`);
export const acceptFriendRequest = (requestId) => api.post(`/friend-requests/${requestId}/accept`);
export const rejectFriendRequest = (requestId) => api.post(`/friend-requests/${requestId}/reject`);
export const getMessages = (username) => api.get(`/messages/${username}`);
export const sendMessage = (messageData) => api.post(`/messages`, messageData);
export const getLists = () => api.get(`/lists`);
export const createList = (listData) => api.post(`/lists`, listData);
export const addToList = (listId, itemData) => api.post(`/lists/${listId}/add`, itemData);
export const getHangouts = () => api.get(`/hangouts`);
export const createHangout = (hangoutData) => api.post(`/hangouts`, hangoutData);
export const joinHangout = (hangoutId) => api.post(`/hangouts/${hangoutId}/join`);
export const leaveHangout = (hangoutId) => api.post(`/hangouts/${hangoutId}/leave`);
export const getPages = () => api.get(`/pages`);
export const createPage = (pageData) => api.post(`/pages`, pageData);
export const followPage = (pageId) => api.post(`/pages/${pageId}/follow`);
export const unfollowPage = (pageId) => api.post(`/pages/${pageId}/unfollow`);

export default api;