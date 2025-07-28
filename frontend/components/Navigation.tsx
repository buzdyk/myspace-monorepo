'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'

interface NavigationProps {
  active: string | null
  children: React.ReactNode
}

export default function Navigation({ active, children }: NavigationProps) {
  const pages = [
    { path: 'today', caption: 'Today' },
    { path: 'month', caption: 'Month' },
    { path: 'settings', caption: 'Settings' },
  ]

  return (
    <div className="absolute w-full" style={{ bottom: '40px' }}>
      {children}

      <div className="mt-4 flex justify-around items-center">
        <div className="flex justify-start">
          {pages.map((page, index) => (
            <div 
              key={page.path}
              className={`text-xs ${index !== pages.length - 1 ? 'mr-6' : ''}`}
            >
              <Link 
                href={`/${page.path}`} 
                className={`hover:text-gray-100 ${
                  page.path === active ? 'text-gray-200' : 'text-gray-600'
                }`}
              >
                {page.caption}
              </Link>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}