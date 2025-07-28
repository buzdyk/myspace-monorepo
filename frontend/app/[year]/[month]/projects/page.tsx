'use client'

import { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import api from '@/lib/api'

interface Project {
  source: string
  project_id: string
  project_title: string
  seconds: number
  hours: number
}

interface ProjectsData {
  year: number
  month: number
  projects: Project[]
  total_hours: number
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

  return (
    <div className="max-w-6xl mx-auto">
      <h1 className="text-3xl font-bold mb-8">Projects - {monthName}</h1>

      <div className="bg-white rounded-lg shadow p-6 mb-8">
        <h2 className="text-xl font-semibold mb-4">Monthly Overview</h2>
        <div className="text-3xl font-bold text-blue-600">
          {data.total_hours.toFixed(2)} total hours
        </div>
      </div>

      <div className="bg-white rounded-lg shadow">
        <div className="p-6">
          <h2 className="text-xl font-semibold mb-4">Project Breakdown</h2>
          
          {data.projects.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              No projects found for this month
            </div>
          ) : (
            <div className="space-y-4">
              {data.projects.map((project, index) => (
                <div key={index} className="border rounded-lg p-4">
                  <div className="flex justify-between items-start">
                    <div>
                      <h3 className="font-semibold text-lg">{project.project_title}</h3>
                      <p className="text-sm text-gray-600">
                        Source: {project.source} | ID: {project.project_id}
                      </p>
                    </div>
                    <div className="text-right">
                      <div className="text-2xl font-bold text-blue-600">
                        {project.hours.toFixed(2)} hrs
                      </div>
                      <div className="text-sm text-gray-500">
                        {project.seconds} seconds
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}