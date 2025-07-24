import './globals.css'
import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import { QueryProviders } from '@/lib/providers'
import { ToastProvider } from '@/components/ui/toast-provider'
import AuthDebugger from '@/components/AuthDebugger'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Zplus SaaS - Multi-tenant Platform',
  description: 'Comprehensive SaaS platform with CRM, HRM, POS, LMS and more modules',
  keywords: 'SaaS, CRM, HRM, POS, LMS, multi-tenant, business platform',
}

interface RootLayoutProps {
  children: React.ReactNode
}

export default function RootLayout({ children }: RootLayoutProps) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <QueryProviders>
          <ToastProvider>
            {children}
            <AuthDebugger />
          </ToastProvider>
        </QueryProviders>
      </body>
    </html>
  )
}
