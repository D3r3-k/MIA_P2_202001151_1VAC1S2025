import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export function middleware(request: NextRequest) {
  const isLoggedIn = request.cookies.get("authToken")?.value === "true";
  const { pathname } = request.nextUrl;
  if (isLoggedIn && pathname === "/login") {
    return NextResponse.redirect(new URL("/", request.url));
  }
  if (!isLoggedIn && pathname.startsWith("/drives/") && pathname !== "/login") {
    return NextResponse.redirect(new URL("/login", request.url));
  }
  return NextResponse.next();
}

export const config = {
  matcher: [
    "/login",
    "/drives/:driveletter*",
    "/drives/:driveletter/:partition_id*",
  ],
};
