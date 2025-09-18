// Auth utility functions for httpOnly cookie-based authentication

/**
 * Make an authenticated API request
 * @param {string} url - The API endpoint
 * @param {object} options - Fetch options
 * @returns {Promise<Response>} - Fetch response
 */
export async function authenticatedFetch(url, options = {}) {
  const defaultOptions = {
    credentials: 'include', // Include cookies in requests
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  };

  const response = await fetch(url, defaultOptions);

  // If we get a 401, the user needs to log in again
  if (response.status === 401) {
    // Redirect to login page
    if (typeof window !== 'undefined') {
      window.location.href = '/login';
    }
    throw new Error('Authentication required');
  }

  return response;
}

/**
 * Login with email and password
 * @param {string} email
 * @param {string} password
 * @returns {Promise<object>} - User data
 */
export async function login(email, password) {
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/auth/login`, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password }),
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || 'Login failed');
  }

  return data;
}

/**
 * Logout the current user
 * @returns {Promise<void>}
 */
export async function logout() {
  try {
    await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/auth/logout`, {
      method: 'POST',
      credentials: 'include',
    });
  } catch (error) {
    // Even if logout fails on server, we should still redirect
    console.error('Logout error:', error);
  }

  // Redirect to login page
  if (typeof window !== 'undefined') {
      window.location.href = '/login';
    }
}

/**
 * Check if user is authenticated by making a test request
 * @returns {Promise<boolean>}
 */
export async function isAuthenticated() {
  try {
    // Make a request to a protected endpoint to verify authentication
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/posts?page=1&page_size=1`, {
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    return response.ok && response.status !== 401;
  } catch {
    return false;
  }
}

/**
 * Register a new user
 * @param {object} userData - User registration data
 * @returns {Promise<object>} - Registration response
 */
export async function register(userData) {
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/auth/register`, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || 'Registration failed');
  }

  return data;
}