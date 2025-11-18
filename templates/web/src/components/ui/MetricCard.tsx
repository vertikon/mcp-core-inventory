import React from 'react'
import { Metric } from '../../types'

interface MetricCardProps {
  metric: Metric
  icon: string
  bgColor: string
  borderColor: string
  iconBgColor: string
  textColor: string
}

export const MetricCard: React.FC<MetricCardProps> = ({
  metric,
  icon,
  bgColor,
  borderColor,
  iconBgColor,
  textColor,
}) => {
  return (
    <div className={`${bgColor} rounded-xl p-4 border ${borderColor} hover:shadow-md transition-all`}>
      <div className="flex items-start justify-between mb-3">
        <div className={`w-10 h-10 ${iconBgColor} rounded-lg flex items-center justify-center border ${borderColor}`}>
          <i className={`${icon} text-xl ${textColor}`}></i>
        </div>
        <div className="text-xs text-gray-500 bg-white px-2 py-1 rounded-md border border-gray-200">
          Meta: {metric.target}
        </div>
      </div>
      <div className="space-y-1">
        <div className="text-xs font-medium text-gray-600 uppercase tracking-wide">
          {metric.name}
        </div>
        <div className="flex items-baseline gap-1">
          <span className={`text-3xl font-bold ${textColor}`}>
            {metric.value.toFixed(1)}
          </span>
          <span className="text-sm font-medium text-gray-500">{metric.unit}</span>
        </div>
      </div>
    </div>
  )
}

