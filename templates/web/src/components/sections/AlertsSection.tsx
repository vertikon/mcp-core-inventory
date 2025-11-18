import React from 'react'
import { Alert } from '../../types'

interface AlertsSectionProps {
  alerts: Alert[]
}

export const AlertsSection: React.FC<AlertsSectionProps> = ({ alerts }) => {
  const hasAlerts = alerts.length > 0

  return (
    <div className="bg-white rounded-2xl shadow-lg p-6 border border-gray-200">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-bold text-gray-800 flex items-center gap-2">
          <i className="ri-alarm-warning-line text-orange-600"></i>
          ALERTAS DO SISTEMA
        </h2>
      </div>
      {hasAlerts ? (
        <div className="space-y-2">
          {alerts.map((alert) => (
            <div
              key={alert.id}
              className={`p-3 rounded-lg border ${
                alert.severity === 'critical'
                  ? 'bg-red-50 border-red-200'
                  : alert.severity === 'warning'
                  ? 'bg-yellow-50 border-yellow-200'
                  : 'bg-blue-50 border-blue-200'
              }`}
            >
              <div className="flex items-center gap-2">
                <i
                  className={`${
                    alert.severity === 'critical'
                      ? 'ri-error-warning-fill text-red-600'
                      : alert.severity === 'warning'
                      ? 'ri-alert-line text-yellow-600'
                      : 'ri-information-line text-blue-600'
                  }`}
                ></i>
                <span className="text-sm font-medium text-gray-800">{alert.message}</span>
                <span className="text-xs text-gray-500 ml-auto">{alert.timestamp}</span>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center py-8">
          <div className="w-16 h-16 bg-green-50 rounded-full flex items-center justify-center mx-auto mb-3 border border-green-200">
            <i className="ri-shield-check-line text-green-600 text-2xl"></i>
          </div>
          <p className="text-green-700 font-medium">🟢 Sem Alertas</p>
          <p className="text-sm text-gray-600 mt-1">Todos os sistemas operando normalmente</p>
        </div>
      )}
    </div>
  )
}

