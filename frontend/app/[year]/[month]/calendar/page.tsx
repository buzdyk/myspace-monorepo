'use client'

import { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import api from '@/lib/api'

interface CalendarDay {
  day: number | null
  hours: number | null
}

interface CalendarData {
  year: number
  month: number
  days: CalendarDay[]
}

export default function CalendarPage() {
  const params = useParams()
  const [data, setData] = useState<CalendarData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await api.get(`/${params.year}/${params.month}/calendar`)
        setData(response.data)
      } catch (err) {
        setError('Failed to fetch calendar data')
        console.error(err)
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [params.year, params.month])

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="text-lg">Loading...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="text-lg text-red-600">{error}</div>
      </div>
    )
  }

  if (!data) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="text-lg">No data available</div>
      </div>
    )
  }

  const monthName = new Date(data.year, data.month - 1).toLocaleDateString('en-US', {
    month: 'long',
    year: 'numeric'
  })

  const weekdays = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
  
  const weeks = []
  for (let i = 0; i < data.days.length; i += 7) {
    weeks.push(data.days.slice(i, i + 7))
  }

  return (
    <div className="max-w-6xl mx-auto">
      <h1 className="text-3xl font-bold mb-8">Calendar - {monthName}</h1>

      <div className="bg-white rounded-lg shadow p-6">
        <div className="grid grid-cols-7 gap-1">
          {weekdays.map(day => (
            <div key={day} className="p-3 text-center font-semibold text-gray-600 border-b">
              {day}
            </div>
          ))}
          
          {data.days.map((day, index) => (
            <div 
              key={index} 
              className={`p-3 min-h-[80px] border border-gray-200 ${
                day.day ? 'bg-white hover:bg-gray-50' : 'bg-gray-100'
              }`}
            >
              {day.day && (
                <>
                  <div className="font-medium text-gray-800">{day.day}</div>
                  {day.hours !== null && (
                    <div className="text-sm text-blue-600 mt-1">
                      {day.hours.toFixed(1)}h
                    </div>
                  )}
                </>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}