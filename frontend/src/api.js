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
export const getFeed = () => api.get('/feed');
export const getExplore = () => api.get('/explore');
export const getUserPhotos = (username) => api.get(`/users/${username}/photos`);
export const createPost = (post) => api.post('/posts', post);
export const likePost = (postId) => api.post(`/posts/${postId}/like`);
export const sharePost = (postId) => api.post(`/posts/${postId}/share`);
export const addComment = (postId, comment) => api.post(`/posts/${postId}/comment`, comment);
export const mintPost = (postId) => api.post(`/posts/${postId}/mint`);
export const mintPhoto = (photoId) => api.post(`/photos/${postId}/mint`);
export const searchUsers = (query) => api.get(`/search/users?q=${query}`);
export const createGroup = (group) => api.post('/groups', group);
export const getGroups = () => api.get('/groups');
export const getGroup = (groupId) => api.get(`/groups/${groupId}`);
export const joinGroup = (groupId) => api.post(`/groups/${groupId}/join`);
export const leaveGroup = (groupId) => api.post(`/groups/${groupId}/leave`);
export const createPage = (page) => api.post('/pages', page);
export const getPages = () => api.get('/pages');
export const getPage = (pageId) => api.get(`/pages/${pageId}`);
export const followPage = (pageId) => api.post(`/pages/${pageId}/follow`);
export const unfollowPage = (pageId) => api.post(`/pages/${pageId}/unfollow`);
export const sendFriendRequest = (request) => api.post('/friend-requests', request);
export const getFriendRequests = () => api.get('/friend-requests');
export const acceptFriendRequest = (requestId) => api.post(`/friend-requests/${requestId}/accept`);
export const rejectFriendRequest = (requestId) => api.post(`/friend-requests/${requestId}/reject`);
export const followUser = (username) => api.post(`/users/${username}/follow`);
export const unfollowUser = (username) => api.post(`/users/${username}/unfollow`);
export const blockUser = (username) => api.post(`/users/${username}/block`);
export const unblockUser = (username) => api.post(`/users/${username}/unblock`);
export const getNotifications = () => api.get('/notifications');
export const markNotificationRead = (notifId) => api.post(`/notifications/${notifId}/read`);
export const sendMessage = (message) => api.post('/messages', message);
export const getMessages = (username) => api.get(`/messages/${username}`);
export const createList = (list) => api.post('/lists', list);
export const getLists = () => api.get('/lists');
export const addToList = (listId, item) => api.post(`/lists/${listId}/add`, item);
export const createHangout = (hangout) => api.post('/hangouts', hangout);
export const getHangouts = () => api.get('/hangouts');
export const joinHangout = (hangoutId) => api.post(`/hangouts/${hangoutId}/join`);
export const leaveHangout = (hangoutId) => api.post(`/hangouts/${hangoutId}/leave`);

export default api;