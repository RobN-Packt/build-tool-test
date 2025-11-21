import { NextRequest } from 'next/server';
import { proxyRequest } from '../proxy';

type RouteContext = {
  params: {
    id: string;
  };
};

function buildPath(id: string) {
  return `/books/${encodeURIComponent(id)}`;
}

export async function GET(request: NextRequest, context: RouteContext) {
  return proxyRequest(request, buildPath(context.params.id));
}

export async function PUT(request: NextRequest, context: RouteContext) {
  return proxyRequest(request, buildPath(context.params.id));
}

export async function DELETE(request: NextRequest, context: RouteContext) {
  return proxyRequest(request, buildPath(context.params.id));
}
