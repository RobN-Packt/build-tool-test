import { NextRequest } from 'next/server';
import { proxyToBooksBackend } from './proxy';

export const runtime = 'nodejs';

// Proxies /api/books requests to the Fargate backend so the browser only talks to the Vercel origin.
export async function GET(req: NextRequest) {
  return proxyToBooksBackend(req, '/books');
}

export async function POST(req: NextRequest) {
  return proxyToBooksBackend(req, '/books');
}
