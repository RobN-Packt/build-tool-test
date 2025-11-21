import { NextRequest } from 'next/server';
import { proxyRequest } from './proxy';

export async function GET(request: NextRequest) {
  return proxyRequest(request, '/books');
}

export async function POST(request: NextRequest) {
  return proxyRequest(request, '/books');
}
