import type { Metadata } from 'next';
import Link from 'next/link';
import './globals.css';

export const metadata: Metadata = {
  title: 'Packt Library',
  description: 'Manage your Packt book inventory with ease.'
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="packt-body">
        <header className="packt-header">
          <div className="container header-content">
            <Link href="/" className="brand">
              Packt Library
            </Link>
            <nav className="packt-nav">
              <Link href="/" className="packt-nav-link">
                Home
              </Link>
              <Link href="/admin/new" className="packt-nav-link accent">
                Add Book
              </Link>
            </nav>
          </div>
        </header>
        <main className="packt-main">
          <div className="container">{children}</div>
        </main>
        <footer className="packt-footer">
          <div className="container packt-footer-content">
            <span>© {new Date().getFullYear()} Packt Publishing</span>
            <span className="text-muted">Building knowledge through world-class technical content.</span>
          </div>
        </footer>
      </body>
    </html>
  );
}
