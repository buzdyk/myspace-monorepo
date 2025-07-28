'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'

export default function TodayRedirect() {
  const router = useRouter()
  
  useEffect(() => {
    const now = new Date()
    const year = now.getFullYear()
    const month = now.getMonth() + 1
    const day = now.getDate()
    
    router.push(`/${year}/${month}/${day}`)
  }, [router])

  return (
    <div className="flex justify-center items-center h-64">
      <div className="text-lg">Redirecting to today...</div>
    </div>
  )
}