'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'

export default function MonthRedirect() {
  const router = useRouter()
  
  useEffect(() => {
    const now = new Date()
    const year = now.getFullYear()
    const month = now.getMonth() + 1
    
    router.push(`/${year}/${month}/projects`)
  }, [router])

  return (
    <div className="flex justify-center items-center h-64">
      <div className="text-lg">Redirecting to current month...</div>
    </div>
  )
}