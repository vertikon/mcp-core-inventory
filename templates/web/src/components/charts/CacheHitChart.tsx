import React from 'react'
import { CacheStats } from '../../types'

interface CacheHitChartProps {
  stats: CacheStats
}

export const CacheHitChart: React.FC<CacheHitChartProps> = ({ stats }) => {
  const levels = [
    { name: 'Cache Geral', value: stats.general, color: 'from-green-500 to-green-600' },
    { name: 'L1 Cache', value: stats.l1, color: 'from-blue-500 to-blue-600' },
    { name: 'L2 Cache', value: stats.l2, color: 'from-purple-500 to-purple-600' },
    { name: 'L3 Cache', value: stats.l3, color: 'from-orange-500 to-orange-600' },
  ]

  return (
    <div className="space-y-4">
      {levels.map((level, index) => (
        <div key={index}>
          <div className="flex justify-between text-sm font-medium text-gray-700 mb-2">
            <span>{level.name}</span>
            <span>{level.value}%</span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-3">
            <div
              className={`bg-gradient-to-r ${level.color} h-3 rounded-full transition-all duration-500`}
              style={{ width: `${level.value}%` }}
            />
          </div>
        </div>
      ))}
      <div className="mt-4 p-3 bg-gray-50 rounded-lg">
        <div className="text-xs text-gray-600 mb-1">Estatísticas</div>
        <div className="flex justify-between text-sm">
          <span className="text-green-600 font-medium">Hit: {stats.hit}%</span>
          <span className="text-red-600 font-medium">Miss: {stats.miss}%</span>
        </div>
      </div>
    </div>
  )
}

