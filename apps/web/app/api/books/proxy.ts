import { NextRequest, NextResponse } from 'next/server';

const METHODS_WITH_BODY = new Set(['POST', 'PUT', 'PATCH', 'DELETE']);

function missingBackendResponse() {
  return NextResponse.json(
    { error: 'BACKEND_INTERNAL_URL is not configured' },
    { status: 500 }
  );
}

function buildUpstreamUrl(path: string, search: string) {
  const backendBase = process.env.BACKEND_INTERNAL_URL;
  if (!backendBase) {
    return null;
  }

  const normalizedBase = backendBase.endsWith('/')
    ? backendBase.slice(0, -1)
    : backendBase;

  return `${normalizedBase}${path}${search}`;
}

export async function proxyToBooksBackend(req: NextRequest, targetPath: string) {
  const upstreamUrl = buildUpstreamUrl(targetPath, req.nextUrl.search);
  if (!upstreamUrl) {
    return missingBackendResponse();
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
