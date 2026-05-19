import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
  const token = request.cookies.get('token')?.value;
  const protectedPaths = ['/dashboard', '/projects', '/reports', '/settings'];
  const isProtectedPath = protectedPaths.some((p) => request.nextUrl.pathname.startsWith(p));

  if (isProtectedPath && !token) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  const authPaths = ['/login', '/register'];
  if (authPaths.some((p) => request.nextUrl.pathname.startsWith(p)) && token) {
    return NextResponse.redirect(new URL('/dashboard', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/dashboard/:path*', '/projects/:path*', '/reports/:path*', '/settings/:path*', '/login', '/register'],
};
