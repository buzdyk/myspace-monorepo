import { hoursToString, formatMoney } from '@/lib/helpers'

interface Project {
  projectTitle: string
  hours: number
}

interface OverviewProps {
  projects: Project[]
  monthHours: number
  projectedIncome: number
  projectedHours: number
  hourlyRate: number
}

export default function Overview({ projects, monthHours, projectedIncome, projectedHours, hourlyRate }: OverviewProps) {
  return (
    <div>
      {projects.map((project, index) => 
        project.hours > 0 && (
          <div key={index} className="flex justify-between">
            <div className="w-3/4 flex justify-between">
              {project.projectTitle}
              <div className="flex-grow border-b-dots mx-3 mb-1 pr-32"></div>
            </div>
            <div className="text-right">{hoursToString(project.hours)}</div>
            <div className="w-32 text-right">
              {formatMoney((project.hours * hourlyRate))}
            </div>
          </div>
        )
      )}

      <div className="mt-8 flex justify-between">
        <div className="w-3/4"></div>
        <div className="text-right">{hoursToString(monthHours)}</div>
        <div className="w-32 text-right">
          {formatMoney((monthHours * hourlyRate))}
        </div>
      </div>

      <div className="flex justify-between text-gray-600">
        <div className="w-3/4 text-gray-600"></div>
        <div className="text-right">{hoursToString(projectedHours)}</div>
        <div className="w-32 text-right">
          {formatMoney(projectedIncome)}
        </div>
      </div>
    </div>
  )
}