import React from 'react'
import { ComponentStatus } from '../../types'

interface ComponentStatusCardProps {
  component: ComponentStatus
}

export const ComponentStatusCard: React.FC<ComponentStatusCardProps> = ({ component }) => {
  const statusColors = {
    online: {
      bg: 'bg-green-50',
      border: 'border-green-200',
      icon: 'text-green-600',
      text: 'text-green-700',
      checkIcon: 'ri-checkbox-circle-fill text-green-600',
    },
    offline: {
      bg: 'bg-red-50',
      border: 'border-red-200',
      icon: 'text-red-600',
      text: 'text-red-700',
      checkIcon: 'ri-close-circle-fill text-red-600',
    },
    degraded: {
      bg: 'bg-yellow-50',
      border: 'border-yellow-200',
      icon: 'text-yellow-600',
      text: 'text-yellow-700',
      checkIcon: 'ri-error-warning-fill text-yellow-600',
    },
  }

  const colors = statusColors[component.status]

  return (
    <div className={`flex items-center justify-between p-3 bg-gray-50 rounded-xl border border-gray-200 hover:shadow-md transition-all`}>
      <div className="flex items-center gap-3">
        <div className={`w-10 h-10 rounded-lg flex items-center justify-center border ${colors.icon} ${colors.bg} ${colors.border}`}>
          <i className={`${component.icon} text-lg`}></i>
        </div>
        <div>
          <div className="font-medium text-gray-800 text-sm">{component.name}</div>
          <div className="text-xs text-gray-600">Uptime: {component.uptime}%</div>
        </div>
      </div>
      <div className="flex items-center gap-2">
        <i className={colors.checkIcon}></i>
        <span className={`text-xs font-medium uppercase ${colors.text}`}>
          {component.status === 'online' ? 'Online' : component.status === 'offline' ? 'Offline' : 'Degradado'}
        </span>
      </div>
    </div>
  )
}

