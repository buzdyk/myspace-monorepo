import { hoursToString } from '@/lib/helpers'

interface TodayCardProps {
  runningHours: number
  todayPercent: number
  todayHours: number
  isToday: boolean
}

export default function TodayCard({ runningHours, todayPercent, todayHours, isToday }: TodayCardProps) {
  return (
    <div>
      <div className="relative text-gray-600">
        Goal
        {runningHours > 0 && (
          <>
            <div 
              className={`absolute rounded-full ${isToday ? 'bg-red-600' : 'bg-red-900'}`}
              style={{ 
                width: '10px', 
                height: '10px', 
                left: '-45px', 
                top: '12px' 
              }}
            />
            <div 
              className="absolute text-sm px-4 py-2 font-bold text-gray-700 hover:text-gray-200 cursor-none"
              style={{ 
                left: '-135px', 
                top: '0' 
              }}
            >
              {hoursToString(runningHours)}
            </div>
          </>
        )}
      </div>

      <div className="mt-4 group w-12">
        <span className="group-hover:hidden">{todayPercent}%</span>
        <span className="text-gray-800 hidden group-hover:inline-block group-hover:text-gray-200">
          {hoursToString(todayHours)}
        </span>
      </div>
    </div>
  )
}