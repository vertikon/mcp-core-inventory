import React from 'react'
import { ChartDataPoint } from '../../types'

interface LineChartProps {
  data: ChartDataPoint[]
  color: string
  label: string
  unit: string
  height?: number
}

export const LineChart: React.FC<LineChartProps> = ({
  data,
  color,
  label,
  unit,
  height = 160,
}) => {
  if (data.length === 0) return null

  const width = 400
  const padding = 20
  const chartWidth = width - padding * 2
  const chartHeight = height - padding * 2

  const maxY = Math.max(...data.map(d => d.y), 1)
  const minY = Math.min(...data.map(d => d.y), 0)

  const scaleY = (value: number) => {
    if (maxY === minY) return chartHeight / 2
    return chartHeight - ((value - minY) / (maxY - minY)) * chartHeight
  }

  const scaleX = (index: number) => {
    return (index / (data.length - 1)) * chartWidth
  }

  const points = data
    .map((point, index) => `${scaleX(index)},${scaleY(point.y)}`)
    .join(' ')

  const areaPoints = [
    `${scaleX(0)},${chartHeight}`,
    ...data.map((point, index) => `${scaleX(index)},${scaleY(point.y)}`),
    `${scaleX(data.length - 1)},${chartHeight}`,
  ].join(' ')

  const currentValue = data[data.length - 1]?.y || 0

  return (
    <div className="relative h-48 bg-gray-50 rounded-xl p-4 overflow-hidden">
      <svg className="w-full h-full" viewBox={`0 0 ${width} ${height}`}>
        <defs>
          <linearGradient id={`gradient-${color}`} x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" stopColor={color} stopOpacity="0.3" />
            <stop offset="100%" stopColor={color} stopOpacity="0.05" />
          </linearGradient>
        </defs>

        {/* Grid lines */}
        {[0, 0.25, 0.5, 0.75, 1].map((ratio) => (
          <line
            key={ratio}
            x1={padding}
            y1={padding + ratio * chartHeight}
            x2={width - padding}
            y2={padding + ratio * chartHeight}
            stroke="#e5e7eb"
            strokeWidth="1"
          />
        ))}

        {/* Area */}
        <polygon
          fill={`url(#gradient-${color})`}
          points={areaPoints}
        />

        {/* Line */}
        <polyline
          fill="none"
          stroke={color}
          strokeWidth="3"
          strokeLinecap="round"
          strokeLinejoin="round"
          points={points}
        />

        {/* Points */}
        {data.map((point, index) => (
          <circle
            key={index}
            cx={scaleX(index)}
            cy={scaleY(point.y)}
            r="4"
            fill={color}
            className="opacity-80"
          />
        ))}
      </svg>
      <div className="absolute top-4 right-4 bg-white px-3 py-1 rounded-lg border border-gray-200">
        <span className="text-sm font-medium text-gray-600">
          {currentValue.toFixed(1)} {unit}
        </span>
      </div>
    </div>
  )
}

