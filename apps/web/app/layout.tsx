import type { Metadata } from 'next';
import Link from 'next/link';
import './globals.css';

export const metadata: Metadata = {
  title: 'Book Store',
  description: 'Simple book catalog powered by Huma API'
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body
        style={{
          margin: 0,
          fontFamily:
            '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, sans-serif',
          backgroundColor: '#f7f7f8'
        }}
      >
        <header
          style={{
            backgroundColor: '#ffffff',
            borderBottom: '1px solid #e5e5e5',
            padding: '1rem 2rem',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center'
          }}
        >
          <Link href="/" style={{ fontSize: '1.25rem', fontWeight: 600, textDecoration: 'none' }}>
            Book Store
          </Link>
          <nav style={{ display: 'flex', gap: '1rem' }}>
            <Link href="/">Home</Link>
            <Link href="/admin/new">New Book</Link>
          </nav>
        </header>
        <main style={{ padding: '2rem', maxWidth: 960, margin: '0 auto' }}>{children}</main>
      </body>
    </html>
  );
}
