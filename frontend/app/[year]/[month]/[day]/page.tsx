'use client'

import { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import api from '@/lib/api'

interface TodayData {
  date: string
  hours: number
  running_hours: number
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

  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-8">
        {new Date(data.date).toLocaleDateString('en-US', {
          weekday: 'long',
          year: 'numeric',
          month: 'long',
          day: 'numeric'
        })}
      </h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold mb-4">Today's Hours</h2>
          <div className="text-3xl font-bold text-blue-600">
            {data.hours.toFixed(2)} hrs
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold mb-4">Currently Running</h2>
          <div className="text-3xl font-bold text-green-600">
            {data.running_hours.toFixed(2)} hrs
          </div>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-xl font-semibold mb-4">Time Breakdown</h2>
        <div className="space-y-2">
          <div className="flex justify-between">
            <span>Total Tracked Time:</span>
            <span className="font-medium">{data.hours.toFixed(2)} hours</span>
          </div>
          <div className="flex justify-between">
            <span>Currently Active:</span>
            <span className="font-medium">{data.running_hours.toFixed(2)} hours</span>
          </div>
          <div className="flex justify-between font-semibold border-t pt-2">
            <span>Combined Total:</span>
            <span>{(data.hours + data.running_hours).toFixed(2)} hours</span>
          </div>
        </div>
      </div>
    </div>
  )
}