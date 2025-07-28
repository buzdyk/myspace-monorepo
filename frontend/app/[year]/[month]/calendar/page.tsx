'use client'

import { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import Navigation from '@/components/Navigation'
import { hoursToString, formatMoney } from '@/lib/helpers'
import api from '@/lib/api'

interface CalendarDay {
  date: number | null
  hours: number | null
  link?: string
  isToday?: boolean
}

interface CalendarData {
  year: number
  month: number
  days: CalendarDay[]
  hours: number
  links: {
    caption: string
    thisLink: string
    prevLink: string
    nextLink: string
  }
  dailyGoal: number
  hourlyRate: number
}

const daysOfWeek = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']

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
      <div className="h-screen w-full flex justify-center items-center font-mono text-2xl selection:bg-red-700 selection:text-white">
        <div>Loading...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="h-screen w-full flex justify-center items-center font-mono text-2xl selection:bg-red-700 selection:text-white">
        <div>{error}</div>
      </div>
    )
  }

  if (!data) {
    return (
      <div className="h-screen w-full flex justify-center items-center font-mono text-2xl selection:bg-red-700 selection:text-white">
        <div>No data available</div>
      </div>
    )
  }

  return (
    <div className="h-screen w-full flex justify-center items-center font-mono text-2xl selection:bg-red-700 selection:text-white">
      <div className="relative">
        <div className="flex justify-start text-xs -ml-1">
          <a href={`${data.links.thisLink}/projects`} className="block pl-1 pr-3 py-1 mr-4 hover:bg-gray-700 text-center text-gray-400">
            Projects
          </a>
          <a href={`${data.links.thisLink}/calendar`} className="block pl-1 pr-3 py-1 ml-4 bg-gray-700 text-center text-gray-400">
            Calendar
          </a>
        </div>

        <div className="mt-10 mb-20">
          <div className="grid grid-cols-7 gap-y-3 gap-x-3 text-center">
            {daysOfWeek.map((day, index) => (
              <div key={index} className="text-left font-bold text-gray-400 text-sm">
                {day}
              </div>
            ))}

            <div className="col-span-7 border-b border-b-gray-700"></div>

            {data.days.map((day, index) => (
              <div key={index} className="group relative w-16">
                {day && (
                  <div className="text-gray-500 text-left text-xs">
                    <a 
                      href={day.link} 
                      className={day.isToday ? 'border-b border-b-red-600' : ''}
                    >
                      {day.date}
                    </a>
                  </div>
                )}

                {day && day.hours && (
                  <div className="flex w-16 cursor-none justify-start mt-2 text-sm">
                    <div className="group-hover:hidden">{hoursToString(day.hours)}</div>
                    <div className="group-hover:block hidden">{formatMoney(day.hours * data.hourlyRate)}</div>
                  </div>
                )}

                {(!day || !day.hours) && <div className="mt-2">&nbsp;</div>}
              </div>
            ))}
          </div>
        </div>
      </div>

      <Navigation active="month">
        <div className="mb-6 text-base flex justify-around">
          <div className="flex justify-around">
            <span className="block text-gray-400">{data.links.caption}</span>

            <a href={`${data.links.prevLink}/calendar`} className="ml-3 text-gray-600 hover:text-gray-200">&lt;</a>
            <a href={`${data.links.nextLink}/calendar`} className="ml-1 text-gray-600 hover:text-gray-200">&gt;</a>
          </div>
        </div>
      </Navigation>
    </div>
  )
}