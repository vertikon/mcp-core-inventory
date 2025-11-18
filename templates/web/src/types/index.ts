export interface Metric {
  name: string
  value: number
  unit: string
  target: string
  status: 'good' | 'warning' | 'critical'
}

export interface ComponentStatus {
  name: string
  status: 'online' | 'offline' | 'degraded'
  uptime: number
  icon: string
}

export interface SystemMetrics {
  throughput: number
  httpP95: number
  cacheHitRatio: number
  circuitBreaker: number
  bootstrapTime: number
  conformity: number
  lastUpdate: string
}

export interface ComponentDetails {
  status: string
  cpuUsage: number
  queueTasks: number
  completionRate: number
  activeThreads: number
  avgTime: number
}

export interface Alert {
  id: string
  severity: 'info' | 'warning' | 'critical'
  message: string
  timestamp: string
}

export interface ChartDataPoint {
  x: number
  y: number
}

export interface CacheStats {
  general: number
  l1: number
  l2: number
  l3: number
  hit: number
  miss: number
}

