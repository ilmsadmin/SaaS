import './globals.css'
import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import { QueryProviders } from '@/lib/providers'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Zplus SaaS - Multi-tenant Platform',
  description: 'Comprehensive SaaS platform with CRM, HRM, POS, LMS and more modules',
  keywords: 'SaaS, CRM, HRM, POS, LMS, multi-tenant, business platform',
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
        </QueryProviders>
      </body>
    </html>
  )
}
