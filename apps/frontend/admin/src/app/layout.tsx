import './globals.css'
import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Toaster } from 'react-hot-toast'
import { AdminProvider } from '@/contexts/AdminContext'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Zplus Admin Dashboard',
  description: 'Admin dashboard for Zplus SaaS Platform',
}

const queryClient = new QueryClient()

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <QueryClientProvider client={queryClient}>
          <AdminProvider>
            {children}
            <Toaster position="top-right" />
          </AdminProvider>
        </QueryClientProvider>
      </body>
    </html>
  )
}
