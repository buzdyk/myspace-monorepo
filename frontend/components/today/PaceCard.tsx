import { hoursToString } from '@/lib/helpers'

interface PaceCardProps {
  pace: number
  dailyGoal: number
}

export default function PaceCard({ pace, dailyGoal }: PaceCardProps) {
  const paceClass = () => {
    if (pace < -dailyGoal) return 'text-red-600'
    if (pace > 0) return 'text-green-600'
    return ''
  }

  return (
    <div>
      <div className="text-gray-600">Pace</div>
      <div className={`mt-4 ${paceClass()}`}>
        {hoursToString(Math.abs(pace))}
      </div>
    </div>
  )
}