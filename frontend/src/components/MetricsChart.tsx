import { useEffect, useRef, useMemo } from 'react';
import ReactECharts from 'echarts-for-react';
import { useContainerStore } from '../stores/containerStore';

interface MetricsChartProps {
  containerId: string;
  showFullTime?: boolean; // If true, show full time format in detail view
}

export default function MetricsChart({ containerId, showFullTime = false }: MetricsChartProps) {
  const metrics = useContainerStore((state) => state.metrics[containerId] || []);
  const chartRef = useRef<ReactECharts>(null);

  // Process metrics - keep last 10 points
  const processedData = useMemo(() => {
    if (metrics.length === 0) return { cpuData: [], memoryData: [] };

    const now = Date.now();
    // Keep last 10 seconds of data
    const filtered = metrics
      .filter((m) => {
        const timestamp = new Date(m.timestamp).getTime();
        return timestamp >= now - 10000; // last 10 seconds
      })
      .slice(-10);

    return {
      cpuData: filtered.map((m) => [m.timestamp, m.cpu_percent]),
      memoryData: filtered.map((m) => [m.timestamp, m.memory_percent]),
    };
  }, [metrics]);

  // Format time label - show full time in detail view, seconds only in card view
  const timeFormatter = useMemo(() => {
    if (showFullTime) {
      return (value: number) => {
        const date = new Date(value);
        const hours = date.getHours().toString().padStart(2, '0');
        const minutes = date.getMinutes().toString().padStart(2, '0');
        const seconds = date.getSeconds().toString().padStart(2, '0');
        return `${hours}:${minutes}:${seconds}`;
      };
    } else {
      // seconds only
      return (value: number) => {
        const date = new Date(value);
        return `${date.getSeconds().toString().padStart(2, '0')}`;
      };
    }
  }, [showFullTime]);

  const option = useMemo(() => ({
    backgroundColor: 'transparent',
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true,
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
      },
      backgroundColor: 'rgba(255, 255, 255, 0.95)',
      borderColor: '#3b82f6',
      borderWidth: 1,
      textStyle: {
        color: '#1e40af',
      },
    },
    legend: {
      data: ['CPU %', 'Memory %'],
      textStyle: {
        color: '#000000',
        fontWeight: 'bold',
      },
      top: 10,
    },
    xAxis: {
      type: 'time',
      boundaryGap: false,
      axisLabel: {
        color: '#000000',
        formatter: timeFormatter,
      },
      axisLine: {
        lineStyle: {
          color: '#93c5fd',
        },
      },
    },
    yAxis: {
      type: 'value',
      max: 100,
      axisLabel: {
        color: '#000000',
        formatter: '{value}%',
      },
      axisLine: {
        lineStyle: {
          color: '#93c5fd',
        },
      },
      splitLine: {
        lineStyle: {
          color: '#dbeafe',
          type: 'dashed',
        },
      },
    },
    series: [
      {
        name: 'CPU %',
        type: 'line',
        smooth: true,
        data: processedData.cpuData,
        itemStyle: {
          color: '#f59e0b',
        },
        lineStyle: {
          width: 3,
          color: '#f59e0b',
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(245, 158, 11, 0.4)' },
              { offset: 1, color: 'rgba(245, 158, 11, 0.1)' },
            ],
          },
        },
      },
      {
        name: 'Memory %',
        type: 'line',
        smooth: true,
        data: processedData.memoryData,
        itemStyle: {
          color: '#10b981',
        },
        lineStyle: {
          width: 3,
          color: '#10b981',
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(16, 185, 129, 0.4)' },
              { offset: 1, color: 'rgba(16, 185, 129, 0.1)' },
            ],
          },
        },
      },
    ],
  }), [processedData, timeFormatter]);

  // Update chart with incremental data
  useEffect(() => {
    if (chartRef.current && metrics.length > 0) {
      const echartsInstance = chartRef.current.getEchartsInstance();
      echartsInstance.setOption(option, { notMerge: false });
    }
  }, [metrics.length, option]);

  if (metrics.length === 0) {
    return (
      <div className="h-64 flex items-center justify-center text-blue-500 font-medium">
        Waiting for metrics...
      </div>
    );
  }

  return (
    <div className="mt-4">
      <ReactECharts
        ref={chartRef}
        option={option}
        style={{ height: '300px', width: '100%' }}
        opts={{ renderer: 'canvas' }}
      />
    </div>
  );
}

