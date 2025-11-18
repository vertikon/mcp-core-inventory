import { useState, useEffect, useCallback } from 'react'
import { SystemMetrics, ComponentStatus, ComponentDetails, Alert, CacheStats } from '../types'
import axios from 'axios'

const API_BASE_URL = '/api/v1'

export const useMetrics = () => {
  const [metrics, setMetrics] = useState<SystemMetrics | null>(null)
  const [components, setComponents] = useState<ComponentStatus[]>([])
  const [componentDetails, setComponentDetails] = useState<ComponentDetails | null>(null)
  const [alerts, setAlerts] = useState<Alert[]>([])
  const [cacheStats, setCacheStats] = useState<CacheStats | null>(null)
  const [isPaused, setIsPaused] = useState(false)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchMetrics = useCallback(async () => {
    if (isPaused) return

    try {
      const [metricsRes, componentsRes, alertsRes, cacheRes] = await Promise.all([
        axios.get(`${API_BASE_URL}/metrics`).catch(() => ({ data: null })),
        axios.get(`${API_BASE_URL}/components/status`).catch(() => ({ data: [] })),
        axios.get(`${API_BASE_URL}/alerts`).catch(() => ({ data: [] })),
        axios.get(`${API_BASE_URL}/cache/stats`).catch(() => ({ data: null })),
      ])

      if (metricsRes.data) {
        setMetrics(metricsRes.data)
      } else {
        // Dados mockados para desenvolvimento
        setMetrics({
          throughput: 445.7,
          httpP95: 41.8,
          cacheHitRatio: 84.6,
          circuitBreaker: 1.2,
          bootstrapTime: 3.1,
          conformity: 98.9,
          lastUpdate: new Date().toLocaleTimeString('pt-BR'),
        })
      }

      if (componentsRes.data && componentsRes.data.length > 0) {
        setComponents(componentsRes.data)
      } else {
        // Dados mockados
        setComponents([
          { name: 'Execution Engine', status: 'online', uptime: 99.8, icon: 'ri-cpu-line' },
          { name: 'Worker Pool', status: 'online', uptime: 99.8, icon: 'ri-team-line' },
          { name: 'Multi-level Cache', status: 'online', uptime: 99.8, icon: 'ri-database-2-line' },
          { name: 'Circuit Breaker', status: 'online', uptime: 99.7, icon: 'ri-shield-check-line' },
          { name: 'Configuration System', status: 'online', uptime: 99.6, icon: 'ri-settings-3-line' },
          { name: 'GLM-4.6 Transformer', status: 'online', uptime: 99.6, icon: 'ri-brain-line' },
          { name: 'Crush Optimizations', status: 'online', uptime: 99.6, icon: 'ri-speed-up-line' },
          { name: 'Observability Stack', status: 'online', uptime: 99.5, icon: 'ri-eye-line' },
        ])
      }

      if (alertsRes.data) {
        setAlerts(alertsRes.data)
      }

      if (cacheRes.data) {
        setCacheStats(cacheRes.data)
      } else {
        setCacheStats({
          general: 87,
          l1: 95,
          l2: 74,
          l3: 62,
          hit: 87,
          miss: 13,
        })
      }

      setError(null)
    } catch (err) {
      setError('Erro ao carregar métricas')
      console.error('Error fetching metrics:', err)
    } finally {
      setLoading(false)
    }
  }, [isPaused])

  const fetchComponentDetails = useCallback(async (componentName: string) => {
    try {
      const res = await axios.get(`${API_BASE_URL}/components/${componentName}/details`)
      if (res.data) {
        setComponentDetails(res.data)
      } else {
        // Dados mockados
        setComponentDetails({
          status: 'Online',
          cpuUsage: 67,
          queueTasks: 156,
          completionRate: 99.8,
          activeThreads: 24,
          avgTime: 23,
        })
      }
    } catch (err) {
      console.error('Error fetching component details:', err)
      setComponentDetails({
        status: 'Online',
        cpuUsage: 67,
        queueTasks: 156,
        completionRate: 99.8,
        activeThreads: 24,
        avgTime: 23,
      })
    }
  }, [])

  useEffect(() => {
    fetchMetrics()
    const interval = setInterval(fetchMetrics, 5000)
    return () => clearInterval(interval)
  }, [fetchMetrics])

  const togglePause = () => {
    setIsPaused(!isPaused)
  }

  return {
    metrics,
    components,
    componentDetails,
    alerts,
    cacheStats,
    isPaused,
    loading,
    error,
    togglePause,
    fetchComponentDetails,
    refresh: fetchMetrics,
  }
}

