import { NextRequest } from 'next/server';

const upstreamBaseUrl =
  process.env.API_SERVER_BASE_URL?.replace(/\/$/, '') ?? 'http://localhost:8080';

const hopByHopHeaders = new Set([
  'connection',
  'keep-alive',
  'proxy-authenticate',
  'proxy-authorization',
  'te',
  'trailer',
  'transfer-encoding',
  'upgrade'
]);

function buildTargetUrl(request: NextRequest, path: string) {
  const { search } = new URL(request.url);
  return `${upstreamBaseUrl}${path}${search}`;
}

function buildForwardHeaders(request: NextRequest) {
  const headers = new Headers();
  request.headers.forEach((value, key) => {
    const headerName = key.toLowerCase();
    if (headerName === 'host' || hopByHopHeaders.has(headerName)) {
      return;
    }
    headers.set(headerName, value);
  });
  return headers;
}

function buildResponseHeaders(upstreamHeaders: Headers) {
  const headers = new Headers();
  upstreamHeaders.forEach((value, key) => {
    const headerName = key.toLowerCase();
    if (hopByHopHeaders.has(headerName)) {
      return;
    }
    if (headerName === 'set-cookie') {
      headers.append(headerName, value);
      return;
    }
    headers.set(headerName, value);
  });
  return headers;
}

export async function proxyRequest(request: NextRequest, path: string) {
  const targetUrl = buildTargetUrl(request, path);
  const headers = buildForwardHeaders(request);

  let body: string | undefined;
  if (!['GET', 'HEAD'].includes(request.method)) {
    body = await request.text();
  }

  try {
    const upstreamResponse = await fetch(targetUrl, {
      method: request.method,
      headers,
      body,
      cache: 'no-store'
    });

    const responseHeaders = buildResponseHeaders(upstreamResponse.headers);

    return new Response(upstreamResponse.body, {
      status: upstreamResponse.status,
      headers: responseHeaders
    });
  } catch (error) {
    console.error('Book API proxy failed', error);
    return new Response(
      JSON.stringify({ error: 'Failed to reach the upstream service.' }),
      {
        status: 502,
        headers: { 'content-type': 'application/json' }
      }
    );
  }
}
