import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'Book Shop',
  description: 'Simple book inventory for managing catalog & stock'
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        {children}
      </body>
    </html>
  );
}
