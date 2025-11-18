import React from 'react'
import { MetricCard } from '../ui/MetricCard'
import { SystemMetrics } from '../../types'

interface MetricsSectionProps {
  metrics: SystemMetrics
}

export const MetricsSection: React.FC<MetricsSectionProps> = ({ metrics }) => {
  const metricCards = [
    {
      metric: {
        name: 'Throughput',
        value: metrics.throughput,
        unit: 'msgs/s',
        target: '200-600',
        status: metrics.throughput >= 200 && metrics.throughput <= 600 ? 'good' : 'warning' as const,
      },
      icon: 'ri-speed-up-line',
      bgColor: 'bg-blue-50',
      borderColor: 'border-blue-200',
      iconBgColor: 'bg-blue-50',
      textColor: 'text-blue-600',
    },
    {
      metric: {
        name: 'HTTP P95',
        value: metrics.httpP95,
        unit: 'ms',
        target: '< 60ms',
        status: metrics.httpP95 < 60 ? 'good' : 'warning' as const,
      },
      icon: 'ri-timer-flash-line',
      bgColor: 'bg-green-50',
      borderColor: 'border-green-200',
      iconBgColor: 'bg-green-50',
      textColor: 'text-green-600',
    },
    {
      metric: {
        name: 'Cache Hit Ratio',
        value: metrics.cacheHitRatio,
        unit: '%',
        target: '70-90%',
        status: metrics.cacheHitRatio >= 70 && metrics.cacheHitRatio <= 90 ? 'good' : 'warning' as const,
      },
      icon: 'ri-database-2-line',
      bgColor: 'bg-purple-50',
      borderColor: 'border-purple-200',
      iconBgColor: 'bg-purple-50',
      textColor: 'text-purple-600',
    },
    {
      metric: {
        name: 'Circuit Breaker',
        value: metrics.circuitBreaker,
        unit: 's',
        target: '< 2s',
        status: metrics.circuitBreaker < 2 ? 'good' : 'warning' as const,
      },
      icon: 'ri-shield-check-line',
      bgColor: 'bg-orange-50',
      borderColor: 'border-orange-200',
      iconBgColor: 'bg-orange-50',
      textColor: 'text-orange-600',
    },
    {
      metric: {
        name: 'Bootstrap Time',
        value: metrics.bootstrapTime,
        unit: 's',
        target: '< 4s',
        status: metrics.bootstrapTime < 4 ? 'good' : 'warning' as const,
      },
      icon: 'ri-rocket-line',
      bgColor: 'bg-teal-50',
      borderColor: 'border-teal-200',
      iconBgColor: 'bg-teal-50',
      textColor: 'text-teal-600',
    },
  ]

  return (
    <div className="bg-white rounded-2xl shadow-lg p-6 border border-gray-200">
      <h2 className="text-lg font-bold text-gray-800 mb-4 flex items-center gap-2">
        <i className="ri-bar-chart-box-line text-blue-600"></i>
        MÉTRICAS PRINCIPAIS
      </h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {metricCards.map((card, index) => (
          <MetricCard key={index} {...card} />
        ))}
      </div>
    </div>
  )
}

