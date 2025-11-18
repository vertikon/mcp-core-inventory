import React from 'react'
import { ComponentStatusCard } from '../ui/ComponentStatusCard'
import { ComponentStatus } from '../../types'

interface ComponentStatusSectionProps {
  components: ComponentStatus[]
}

export const ComponentStatusSection: React.FC<ComponentStatusSectionProps> = ({ components }) => {
  return (
    <div className="bg-white rounded-2xl shadow-lg p-6 border border-gray-200">
      <h2 className="text-lg font-bold text-gray-800 mb-4 flex items-center gap-2">
        <i className="ri-shield-check-line text-green-600"></i>
        STATUS DOS COMPONENTES
      </h2>
      <div className="space-y-3">
        {components.map((component, index) => (
          <ComponentStatusCard key={index} component={component} />
        ))}
      </div>
    </div>
  )
}

