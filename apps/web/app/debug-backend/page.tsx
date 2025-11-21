import { fetchBackendHealthDiagnostics } from '@/lib/api/server';

export const revalidate = 0;

export default async function DebugBackendPage() {
  const diagnostics = await fetchBackendHealthDiagnostics();
  const backendUrl =
    process.env.BACKEND_INTERNAL_URL ??
    process.env.NEXT_PUBLIC_API_BASE_URL ??
    'unset';
  const backendEnvVar = process.env.BACKEND_INTERNAL_URL
    ? 'BACKEND_INTERNAL_URL'
    : process.env.NEXT_PUBLIC_API_BASE_URL
      ? 'NEXT_PUBLIC_API_BASE_URL'
      : 'unset';

  return (
    <div className="page">
      <section className="card">
        <h1>Backend diagnostics</h1>
          <p className="text-muted">
            SSR environment will attempt to reach <code>{diagnostics.url}</code> using{' '}
            <code>{backendEnvVar}</code> = <code>{backendUrl}</code>.
          </p>
        <dl>
          <dt>Status</dt>
          <dd>
            {diagnostics.status} {diagnostics.ok ? '(ok)' : '(error)'}
          </dd>
          <dt>Fetched at</dt>
          <dd>{diagnostics.fetchedAt}</dd>
          {diagnostics.error ? (
            <>
              <dt>Error</dt>
              <dd className="text-danger">{diagnostics.error}</dd>
            </>
          ) : null}
        </dl>
      </section>

      <section className="card">
        <h2>Backend response body</h2>
        <pre>{diagnostics.body || '(no body returned)'}</pre>
      </section>
    </div>
  );
}
