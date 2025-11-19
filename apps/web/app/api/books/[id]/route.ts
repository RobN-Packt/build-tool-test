import { NextRequest } from 'next/server';
import { proxyToBooksBackend } from '../proxy';

type RouteContext = {
  params: {
    id: string;
  };
};

function targetPath(id: string) {
  return `/books/${encodeURIComponent(id)}`;
}

// Proxies /api/books/:id operations to the backend to keep browser traffic on the Vercel domain.
export async function GET(req: NextRequest, context: RouteContext) {
  return proxyToBooksBackend(req, targetPath(context.params.id));
}

export async function PUT(req: NextRequest, context: RouteContext) {
  return proxyToBooksBackend(req, targetPath(context.params.id));
}

export async function DELETE(req: NextRequest, context: RouteContext) {
  return proxyToBooksBackend(req, targetPath(context.params.id));
}
