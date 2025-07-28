'use client'

import { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import Navigation from '@/components/Navigation'
import EmptyPlaceholder from '@/components/projects/EmptyPlaceholder'
import Overview from '@/components/projects/Overview'
import api from '@/lib/api'

interface Project {
  source: string
  project_id: string
  projectTitle: string
  seconds: number
  hours: number
}

interface ProjectsData {
  projects: Project[]
  monthHours: number
  projectedIncome: number
  projectedHours: number
  dayOfMonth: string
  weekdays: number
  weekends: number
  links: {
    caption: string
    thisLink: string
    prevLink: string
    nextLink: string
  }
  hourlyRate: number
}

export default function ProjectsPage() {
  const params = useParams()
  const [data, setData] = useState<ProjectsData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await api.get(`/${params.year}/${params.month}/projects`)
        setData(response.data)
      } catch (err) {
        setError('Failed to fetch projects data')
        console.error(err)
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [params.year, params.month])

  if (loading) {
    return (
      <div className="h-screen flex justify-center items-center font-mono">
        <div className="text-lg selection:bg-red-700 selection:text-white">Loading...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="h-screen flex justify-center items-center font-mono">
        <div className="text-lg selection:bg-red-700 selection:text-white">{error}</div>
      </div>
    )
  }

  if (!data) {
    return (
      <div className="h-screen flex justify-center items-center font-mono">
        <div className="text-lg selection:bg-red-700 selection:text-white">No data available</div>
      </div>
    )
  }

  return (
    <div className="h-screen flex justify-center items-center font-mono">
      <div className="text-lg selection:bg-red-700 selection:text-white">
        <div className="flex justify-start text-xs -ml-1">
          <a href={`${data.links.thisLink}/projects`} className="block pl-1 pr-5 py-1 mr-4 bg-gray-700 text-center text-gray-400">
            Projects
          </a>
          <a href={`${data.links.thisLink}/calendar`} className="block pl-1 pr-5 py-1 ml-4 hover:bg-gray-700 text-center text-gray-400">
            Calendar
          </a>
        </div>

        <div className="mt-10 mb-20">
          {data.monthHours ? (
            <Overview
              projects={data.projects}
              monthHours={data.monthHours}
              projectedIncome={data.projectedIncome}
              projectedHours={data.projectedHours}
              hourlyRate={data.hourlyRate}
            />
          ) : (
            <EmptyPlaceholder />
          )}
        </div>
      </div>

      <Navigation active="month">
        <div className="mb-6 text-base flex justify-around">
          <div className="flex justify-around">
            <span className="block text-gray-400">{data.links.caption}</span>

            <a href={`${data.links.prevLink}/projects`} className="ml-3 text-gray-600 hover:text-gray-200">&lt;</a>
            <a href={`${data.links.nextLink}/projects`} className="ml-1 text-gray-600 hover:text-gray-200">&gt;</a>
          </div>
        </div>
      </Navigation>
    </div>
  )
}