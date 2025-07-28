import type { Metadata } from 'next'
import './globals.css'
import Navigation from '@/components/Navigation'

export const metadata: Metadata = {
  title: 'MySpace - Time Tracking',
  description: 'Unified time tracking from multiple providers',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className="bg-gray-800 text-gray-200">
        {children}
      </body>
    </html>
  )
}