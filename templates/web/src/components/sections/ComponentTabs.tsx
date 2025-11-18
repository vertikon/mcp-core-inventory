import React, { useState } from 'react'
import { ComponentDetails } from '../../types'

interface ComponentTabsProps {
  componentDetails: ComponentDetails | null
  onComponentSelect: (component: string) => void
}

const tabs = [
  { id: 'execution-engine', label: 'EXECUTION ENGINE', icon: 'ri-cpu-line' },
  { id: 'cache', label: 'CACHE', icon: 'ri-database-2-line' },
  { id: 'metrics', label: 'METRICS', icon: 'ri-line-chart-line' },
  { id: 'config', label: 'CONFIG', icon: 'ri-settings-3-line' },
  { id: 'scheduler', label: 'SCHEDULER', icon: 'ri-calendar-line' },
  { id: 'ai', label: 'AI', icon: 'ri-brain-line' },
  { id: 'state', label: 'STATE', icon: 'ri-file-list-3-line' },
]

export const ComponentTabs: React.FC<ComponentTabsProps> = ({
  componentDetails,
  onComponentSelect,
}) => {
  const [activeTab, setActiveTab] = useState('execution-engine')

  const handleTabClick = (tabId: string) => {
    setActiveTab(tabId)
    onComponentSelect(tabId)
  }

  const renderTabContent = () => {
    if (!componentDetails) {
      return (
        <div className="p-6">
          <div className="mb-6">
            <h3 className="text-xl font-bold text-gray-800 mb-2">Execution Engine</h3>
            <p className="text-gray-600 text-sm">Motor de execução principal do sistema</p>
          </div>
          <div className="text-center py-8 text-gray-500">Carregando detalhes...</div>
        </div>
      )
    }

    return (
      <div className="p-6">
        <div className="mb-6">
          <h3 className="text-xl font-bold text-gray-800 mb-2">Execution Engine</h3>
          <p className="text-gray-600 text-sm">Motor de execução principal do sistema</p>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <div className="rounded-xl p-4 border hover:shadow-md transition-all text-green-600 bg-green-50 border-green-200">
            <div className="flex items-center gap-3 mb-2">
              <div className="w-10 h-10 rounded-lg flex items-center justify-center border text-green-600 bg-green-50 border-green-200">
                <i className="ri-play-circle-line text-lg"></i>
              </div>
              <div className="text-xs font-medium uppercase tracking-wide">Status do Motor</div>
            </div>
            <div className="text-2xl font-bold">{componentDetails.status}</div>
          </div>
          <div className="rounded-xl p-4 border hover:shadow-md transition-all text-blue-600 bg-blue-50 border-blue-200">
            <div className="flex items-center gap-3 mb-2">
              <div className="w-10 h-10 rounded-lg flex items-center justify-center border text-blue-600 bg-blue-50 border-blue-200">
                <i className="ri-cpu-line text-lg"></i>
              </div>
              <div className="text-xs font-medium uppercase tracking-wide">Utilização de CPU</div>
            </div>
            <div className="text-2xl font-bold">{componentDetails.cpuUsage}%</div>
          </div>
          <div className="rounded-xl p-4 border hover:shadow-md transition-all text-blue-600 bg-blue-50 border-blue-200">
            <div className="flex items-center gap-3 mb-2">
              <div className="w-10 h-10 rounded-lg flex items-center justify-center border text-blue-600 bg-blue-50 border-blue-200">
                <i className="ri-list-check text-lg"></i>
              </div>
              <div className="text-xs font-medium uppercase tracking-wide">Tarefas em Fila</div>
            </div>
            <div className="text-2xl font-bold">{componentDetails.queueTasks}</div>
          </div>
          <div className="rounded-xl p-4 border hover:shadow-md transition-all text-green-600 bg-green-50 border-green-200">
            <div className="flex items-center gap-3 mb-2">
              <div className="w-10 h-10 rounded-lg flex items-center justify-center border text-green-600 bg-green-50 border-green-200">
                <i className="ri-checkbox-circle-line text-lg"></i>
              </div>
              <div className="text-xs font-medium uppercase tracking-wide">Taxa de Conclusão</div>
            </div>
            <div className="text-2xl font-bold">{componentDetails.completionRate}%</div>
          </div>
          <div className="rounded-xl p-4 border hover:shadow-md transition-all text-blue-600 bg-blue-50 border-blue-200">
            <div className="flex items-center gap-3 mb-2">
              <div className="w-10 h-10 rounded-lg flex items-center justify-center border text-blue-600 bg-blue-50 border-blue-200">
                <i className="ri-play-circle-line text-lg"></i>
              </div>
              <div className="text-xs font-medium uppercase tracking-wide">Threads Ativas</div>
            </div>
            <div className="text-2xl font-bold">{componentDetails.activeThreads}</div>
          </div>
          <div className="rounded-xl p-4 border hover:shadow-md transition-all text-green-600 bg-green-50 border-green-200">
            <div className="flex items-center gap-3 mb-2">
              <div className="w-10 h-10 rounded-lg flex items-center justify-center border text-green-600 bg-green-50 border-green-200">
                <i className="ri-timer-line text-lg"></i>
              </div>
              <div className="text-xs font-medium uppercase tracking-wide">Tempo Médio</div>
            </div>
            <div className="text-2xl font-bold">{componentDetails.avgTime}ms</div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="bg-white rounded-2xl shadow-lg border border-gray-200 overflow-hidden">
      <div className="border-b border-gray-200 bg-gray-50">
        <div className="flex overflow-x-auto scrollbar-hide">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => handleTabClick(tab.id)}
              className={`flex items-center gap-2 px-6 py-4 text-sm font-semibold whitespace-nowrap transition-all border-b-2 ${
                activeTab === tab.id
                  ? 'border-blue-600 text-blue-600 bg-white'
                  : 'border-transparent text-gray-600 hover:text-gray-800 hover:bg-gray-100'
              }`}
            >
              <i className={`${tab.icon} text-lg`}></i>
              {tab.label}
            </button>
          ))}
        </div>
      </div>
      {renderTabContent()}
    </div>
  )
}

