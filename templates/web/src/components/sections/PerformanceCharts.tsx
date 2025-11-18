import React from 'react'
import { LineChart } from '../charts/LineChart'
import { CacheHitChart } from '../charts/CacheHitChart'
import { useChartData } from '../../hooks/useChartData'
import { CacheStats } from '../../types'

interface PerformanceChartsProps {
  cacheStats: CacheStats | null
}

export const PerformanceCharts: React.FC<PerformanceChartsProps> = ({ cacheStats }) => {
  const throughputData = useChartData('throughput')
  const latencyData = useChartData('latency')

  return (
    <div className="space-y-6">
      <h2 className="text-xl font-bold text-gray-800 flex items-center gap-2">
        <i className="ri-line-chart-line text-blue-600"></i>
        GRÁFICOS DE PERFORMANCE
      </h2>
      <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6">
        <div className="bg-white rounded-2xl shadow-lg p-6 border border-gray-200">
          <h3 className="text-lg font-bold text-gray-800 mb-4">Throughput ao Longo do Tempo</h3>
          <LineChart
            data={throughputData}
            color="#3b82f6"
            label="Throughput"
            unit="msgs/s"
          />
          <div className="flex justify-between text-xs text-gray-500 mt-2">
            <span>{new Date(Date.now() - 60000).toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })}</span>
            <span>{new Date().toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })}</span>
          </div>
        </div>
        <div className="bg-white rounded-2xl shadow-lg p-6 border border-gray-200">
          <h3 className="text-lg font-bold text-gray-800 mb-4">Latência HTTP P95 ao Longo do Tempo</h3>
          <LineChart
            data={latencyData}
            color="#10b981"
            label="Latência"
            unit="ms"
          />
          <div className="flex justify-between text-xs text-gray-500 mt-2">
            <span>{new Date(Date.now() - 60000).toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })}</span>
            <span>{new Date().toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })}</span>
          </div>
        </div>
        <div className="bg-white rounded-2xl shadow-lg p-6 border border-gray-200">
          <h3 className="text-lg font-bold text-gray-800 mb-4">Cache Hit Ratio</h3>
          {cacheStats ? (
            <CacheHitChart stats={cacheStats} />
          ) : (
            <div className="text-center py-8 text-gray-500">Carregando dados...</div>
          )}
        </div>
      </div>
    </div>
  )
}

