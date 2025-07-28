'use client'

import { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import Navigation from '@/components/Navigation'
import TodayCard from '@/components/today/TodayCard'
import MonthCard from '@/components/today/MonthCard'
import PaceCard from '@/components/today/PaceCard'
import api from '@/lib/api'

interface TodayData {
  date: string
  hours: number
  running_hours: number
  today_percent: number
  month_percent: number
  month_hours: number
  pace: number
  daily_goal: number
  nav: {
    month: string
    day: string
    year: string
    month_link: string
    prev_link: string
    next_link: string
  }
  is_today: boolean
}

export default function TodayPage() {
  const params = useParams()
  const [data, setData] = useState<TodayData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await api.get(`/${params.year}/${params.month}/${params.day}`)
        setData(response.data)
      } catch (err) {
        setError('Failed to fetch data')
        console.error(err)
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [params.year, params.month, params.day])

  useEffect(() => {
    const timer = setTimeout(() => {
      window.location.reload()
    }, 120000) // 2 minutes like the original

    return () => clearTimeout(timer)
  }, [])

  if (loading) {
    return (
      <div className="h-screen flex justify-around items-center font-mono text-xl selection:bg-red-700 selection:text-white">
        <div>Loading...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="h-screen flex justify-around items-center font-mono text-xl selection:bg-red-700 selection:text-white">
        <div>{error}</div>
      </div>
    )
  }

  if (!data) {
    return (
      <div className="h-screen flex justify-around items-center font-mono text-xl selection:bg-red-700 selection:text-white">
        <div>No data available</div>
      </div>
    )
  }

  return (
    <div className="h-screen flex justify-around items-center font-mono text-xl selection:bg-red-700 selection:text-white">
      <div>
        <div className="flex justify-between w-96">
          <TodayCard
            isToday={data.is_today}
            runningHours={data.running_hours}
            todayHours={data.hours}
            todayPercent={data.today_percent}
          />
          <MonthCard monthPercent={data.month_percent} monthHours={data.month_hours} />
          <PaceCard pace={data.pace} dailyGoal={data.daily_goal} />
        </div>
      </div>

      <Navigation active={data.is_today ? 'today' : null}>
        <div className="mb-6 text-base flex justify-around">
          <div className="relative flex justify-between items-center text-gray-400">
            <a href={data.nav.month_link} className="block">
              {data.nav.month}&nbsp;
            </a> 
            {data.nav.day} {data.nav.year}

            <a href={data.nav.prev_link} className="ml-3 text-gray-600 hover:text-gray-200">&lt;</a>
            <a href={data.nav.next_link} className="ml-1 text-gray-600 hover:text-gray-200">&gt;</a>
          </div>
        </div>
      </Navigation>
    </div>
  )
}