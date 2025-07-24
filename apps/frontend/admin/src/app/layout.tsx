import './globals.css'
import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import { Toaster } from 'react-hot-toast'
import { QueryProviders } from '@/lib/providers'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Zplus Admin Dashboard',
  description: 'Admin dashboard for Zplus SaaS Platform',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <QueryProviders>
          {children}
          <Toaster position="top-right" />
        </QueryProviders>
      </body>
    </html>
  )
}
