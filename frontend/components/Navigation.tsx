'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'

export default function Navigation() {
  const pathname = usePathname()
  
  const isActive = (path: string) => {
    return pathname === path || pathname.startsWith(path)
  }
  
  return (
    <nav className="bg-white shadow-lg">
      <div className="container mx-auto px-4">
        <div className="flex justify-between items-center py-4">
          <Link href="/" className="text-xl font-bold text-gray-800">
            MySpace
          </Link>
          
          <div className="flex space-x-6">
            <Link 
              href="/today" 
              className={`px-3 py-2 rounded-md ${
                isActive('/today') || pathname.match(/\/\d{4}\/\d{1,2}\/\d{1,2}$/)
                  ? 'bg-blue-500 text-white' 
                  : 'text-gray-600 hover:text-gray-800'
              }`}
            >
              Today
            </Link>
            
            <Link 
              href="/month" 
              className={`px-3 py-2 rounded-md ${
                pathname.includes('/projects') || pathname.includes('/calendar')
                  ? 'bg-blue-500 text-white' 
                  : 'text-gray-600 hover:text-gray-800'
              }`}
            >
              Month
            </Link>
          </div>
        </div>
      </div>
    </nav>
  )
}