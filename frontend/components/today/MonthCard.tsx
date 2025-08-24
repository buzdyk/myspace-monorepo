import { hoursToString } from '@/lib/helpers'

interface MonthCardProps {
  monthHours: number
  monthPercent: number
}

export default function MonthCard({ monthHours, monthPercent }: MonthCardProps) {
  return (
    <div>
      <div className="text-gray-600">Month</div>
      <div className="mt-4 text-left group w-12">
        <span className="group-hover:hidden">{monthPercent}%</span>
        <span className="text-gray-800 hidden group-hover:inline-block group-hover:text-gray-200">
          {hoursToString(monthHours)}
        </span>
      </div>
    </div>
  )
}