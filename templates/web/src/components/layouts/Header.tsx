import React from 'react'

interface HeaderProps {
  conformity: number
  lastUpdate: string
  isPaused: boolean
  onPauseToggle: () => void
}

export const Header: React.FC<HeaderProps> = ({
  conformity,
  lastUpdate,
  isPaused,
  onPauseToggle,
}) => {
  return (
    <div className="bg-white rounded-2xl shadow-lg p-6 border border-gray-200">
      <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
        <div className="flex items-center gap-4">
          <div className="w-12 h-12 bg-gradient-to-br from-green-500 to-green-600 rounded-xl flex items-center justify-center">
            <i className="ri-cpu-line text-white text-2xl"></i>
          </div>
          <div>
            <h1 className="text-2xl font-bold text-gray-800">BLOCO-1 CORE PLATFORM</h1>
            <p className="text-sm text-gray-600">Sistema de Monitoramento em Tempo Real</p>
          </div>
        </div>
        <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4">
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2 bg-green-50 px-3 py-2 rounded-lg border border-green-200">
              <div className="w-3 h-3 bg-green-500 rounded-full animate-pulse"></div>
              <span className="text-sm font-semibold text-green-700">
                CONFORMIDADE: {conformity.toFixed(1)}%
              </span>
            </div>
            <div className="flex items-center gap-2 bg-blue-50 px-3 py-2 rounded-lg border border-blue-200">
              <i className="ri-time-line text-blue-600"></i>
              <span className="text-sm font-medium text-blue-700">
                ÚLTIMA ATUALIZAÇÃO: {lastUpdate}
              </span>
            </div>
          </div>
          <button
            onClick={onPauseToggle}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-all ${
              isPaused
                ? 'bg-gray-600 text-white hover:bg-gray-700'
                : 'bg-green-600 text-white hover:bg-green-700'
            }`}
          >
            <i className={isPaused ? 'ri-play-line' : 'ri-pause-line'}></i>
            {isPaused ? 'Retomar' : 'Pausar'}
          </button>
        </div>
      </div>
    </div>
  )
}

