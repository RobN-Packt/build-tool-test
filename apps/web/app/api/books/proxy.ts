import { NextRequest, NextResponse } from 'next/server';
import { getBackendBaseUrl } from '@/lib/api/server';

const METHODS_WITH_BODY = new Set(['POST', 'PUT', 'PATCH', 'DELETE']);

function missingBackendResponse(message?: string) {
  return NextResponse.json(
    { error: message ?? 'BACKEND_INTERNAL_URL is not configured' },
    { status: 500 }
  );
}

function buildUpstreamUrl(path: string, search: string) {
  const backendBase = getBackendBaseUrl();
  return `${backendBase}${path}${search}`;
}

export async function proxyToBooksBackend(req: NextRequest, targetPath: string) {
  let upstreamUrl: string;
  try {
    upstreamUrl = buildUpstreamUrl(targetPath, req.nextUrl.search);
  } catch (error) {
    const message =
      error instanceof Error ? error.message : 'BACKEND_INTERNAL_URL is not configured';
    return missingBackendResponse(message);
  }

  const init: RequestInit = {
    method: req.method,
    cache: 'no-store'
  };

  if (METHODS_WITH_BODY.has(req.method)) {
    init.headers = {
      'Content-Type': req.headers.get('content-type') ?? 'application/json'
    };
    init.body = await req.text();
  }

  try {
    const upstreamResponse = await fetch(upstreamUrl, init);
    const text = await upstreamResponse.text();

    return new NextResponse(text, {
      status: upstreamResponse.status,
      headers: {
        'Content-Type':
          upstreamResponse.headers.get('content-type') ?? 'application/json'
      }
    });
  } catch (error) {
    console.error('Failed to reach backend service', error);
    return NextResponse.json(
      { error: 'Failed to reach backend service' },
      { status: 502 }
    );
  }
}
