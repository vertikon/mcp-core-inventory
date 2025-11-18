import { useState, useEffect } from 'react'
import { ChartDataPoint } from '../types'

const generateRandomData = (count: number, min: number, max: number): ChartDataPoint[] => {
  return Array.from({ length: count }, (_, i) => ({
    x: i,
    y: Math.random() * (max - min) + min,
  }))
}

export const useChartData = (metricType: 'throughput' | 'latency' | 'cache') => {
  const [data, setData] = useState<ChartDataPoint[]>([])

  useEffect(() => {
    let interval: number

    const updateData = () => {
      let newData: ChartDataPoint[]
      
      switch (metricType) {
        case 'throughput':
          newData = generateRandomData(20, 200, 600)
          break
        case 'latency':
          newData = generateRandomData(20, 0, 60)
          break
        case 'cache':
          newData = generateRandomData(20, 70, 95)
          break
        default:
          newData = []
      }
      
      setData(newData)
    }

    updateData()
    interval = setInterval(updateData, 5000)

    return () => {
      if (interval) clearInterval(interval)
    }
  }, [metricType])

  return data
}

