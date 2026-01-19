import { useEffect, useRef } from 'react';
import ReactECharts from 'echarts-for-react';
import { useContainerStore } from '../stores/containerStore';

interface MetricsChartProps {
  containerId: string;
}

export default function MetricsChart({ containerId }: MetricsChartProps) {
  const metrics = useContainerStore((state) => state.metrics[containerId] || []);
  const chartRef = useRef<ReactECharts>(null);

  const cpuData = metrics.map((m) => [m.timestamp, m.cpu_percent]);
  const memoryData = metrics.map((m) => [m.timestamp, m.memory_percent]);

  const option = {
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
    },
    legend: {
      data: ['CPU %', 'Memory %'],
      textStyle: {
        color: '#fff',
      },
    },
    xAxis: {
      type: 'time',
      boundaryGap: false,
      axisLabel: {
        color: '#999',
      },
      axisLine: {
        lineStyle: {
          color: '#333',
        },
      },
    },
    yAxis: {
      type: 'value',
      max: 100,
      axisLabel: {
        color: '#999',
        formatter: '{value}%',
      },
      axisLine: {
        lineStyle: {
          color: '#333',
        },
      },
      splitLine: {
        lineStyle: {
          color: '#333',
        },
      },
    },
    series: [
      {
        name: 'CPU %',
        type: 'line',
        smooth: true,
        data: cpuData,
        itemStyle: {
          color: '#5470c6',
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(84, 112, 198, 0.3)' },
              { offset: 1, color: 'rgba(84, 112, 198, 0.1)' },
            ],
          },
        },
      },
      {
        name: 'Memory %',
        type: 'line',
        smooth: true,
        data: memoryData,
        itemStyle: {
          color: '#91cc75',
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(145, 204, 117, 0.3)' },
              { offset: 1, color: 'rgba(145, 204, 117, 0.1)' },
            ],
          },
        },
      },
    ],
  };

  // Update chart with incremental data
  useEffect(() => {
    if (chartRef.current && metrics.length > 0) {
      const echartsInstance = chartRef.current.getEchartsInstance();
      echartsInstance.setOption(option, { notMerge: false });
    }
  }, [metrics.length]);

  if (metrics.length === 0) {
    return (
      <div className="h-64 flex items-center justify-center text-gray-500">
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

