import { NextResponse } from 'next/server';

export async function middleware(request) {
  const { pathname } = request.nextUrl;

  // Define public routes that don't require authentication
  const publicRoutes = ['/login', '/register'];

  // Define protected routes that require authentication
  const protectedRoutes = ['/posts'];

  // Check if current path is a public route
  const isPublicRoute = publicRoutes.some(route => pathname.startsWith(route));

  // Check if current path is a protected route
  const isProtectedRoute = protectedRoutes.some(route => pathname.startsWith(route));

  // Get the auth token from cookies
  const authToken = request.cookies.get('auth_token')?.value;

  // If accessing a protected route without auth token, redirect to login
  if (isProtectedRoute && !authToken) {
    const loginUrl = new URL('/login', request.url);
    return NextResponse.redirect(loginUrl);
  }

  // If accessing a public route with auth token, redirect to posts
  if (isPublicRoute && authToken) {
    // Verify token is valid by making a quick API call
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/posts?page=1&page_size=1`, {
        headers: {
          'Cookie': `auth_token=${authToken}`,
        },
      });

      // If token is valid, redirect authenticated users away from login/register
      if (response.ok) {
        const postsUrl = new URL('/posts', request.url);
        return NextResponse.redirect(postsUrl);
      }
    } catch (error) {
      // If API call fails, continue to public route
      console.error('Auth verification failed:', error);
    }
  }

  // Handle root path redirects
  if (pathname === '/') {
    if (authToken) {
      // Try to verify token
      try {
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/posts?page=1&page_size=1`, {
          headers: {
            'Cookie': `auth_token=${authToken}`,
          },
        });

        if (response.ok) {
          // Token is valid, redirect to posts
          const postsUrl = new URL('/posts', request.url);
          return NextResponse.redirect(postsUrl);
        } else {
          // Token is invalid, redirect to login
          const loginUrl = new URL('/login', request.url);
          return NextResponse.redirect(loginUrl);
        }
      } catch (error) {
        // API call failed, redirect to login
        const loginUrl = new URL('/login', request.url);
        return NextResponse.redirect(loginUrl);
      }
    } else {
      // No token, redirect to login
      const loginUrl = new URL('/login', request.url);
      return NextResponse.redirect(loginUrl);
    }
  }

  // Allow the request to continue
  return NextResponse.next();
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     */
    '/((?!api|_next/static|_next/image|favicon.ico).*)',
  ],
};