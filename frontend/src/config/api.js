/**
 * API configuration for the Codenames Game
 * 
 * This file defines the base URL for API requests.
 * When running the app with the Go backend, this can be empty
 * to use relative URLs that automatically resolve to the current domain.
 */

// In development or when using a separate API server, you might want to specify the full URL
// const API_URL = 'http://localhost:8080';

// For production, use a relative path when API is served from the same domain as the frontend
const API_URL = 'http://localhost:8080';

export default API_URL;