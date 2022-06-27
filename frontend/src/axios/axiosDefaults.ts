import axios from 'axios';

export const isDev = process.env.NODE_ENV === 'development';

const devUrl = 'http://localhost:4000';
const prodUrl = '';

export const getBaseUrl = () => (isDev ? devUrl : prodUrl);

axios.defaults.baseURL = getBaseUrl();
