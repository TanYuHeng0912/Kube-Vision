import { useState } from 'react';
import { useContainerStore, ContainerMetrics } from '../stores/containerStore';
import { useWebSocket } from '../hooks/useWebSocket';
import { useIntersectionObserver } from '../hooks/useIntersectionObserver';
import MetricsChart from './MetricsChart';
import ContainerDetail from './ContainerDetail';

interface ContainerCardProps {
  container: {
    id: string;
    name: string;
    image: string;
    status: string;
    state: string;
  };
}

export default function ContainerCard({ container }: ContainerCardProps) {
  const { addMetrics } = useContainerStore();
  const [cardRef, isVisible] = useIntersectionObserver<HTMLDivElement>();
  const [showDetail, setShowDetail] = useState(false);

  const wsUrl = isVisible
    ? `ws://localhost:8080/ws/stats/${container.id}`
    : null;

  const { status } = useWebSocket({
    url: wsUrl,
    onMessage: (data: ContainerMetrics) => {
      addMetrics(container.id, data);
    },
    reconnect: true,
  });

  const getStatusColor = () => {
    switch (container.state) {
      case 'running':
        return 'bg-green-500';
      case 'exited':
        return 'bg-red-500';
      case 'paused':
        return 'bg-yellow-500';
      default:
        return 'bg-gray-500';
    }
  };

  const getConnectionStatusColor = () => {
    switch (status) {
      case 'connected':
        return 'bg-green-500';
      case 'connecting':
        return 'bg-yellow-500';
      case 'error':
      case 'disconnected':
        return 'bg-red-500';
      default:
        return 'bg-gray-500';
    }
  };

  return (
    <div
      ref={cardRef}
      className="bg-gray-800 rounded-lg p-6 shadow-lg border border-gray-700"
    >
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-3">
          <div className={`w-3 h-3 rounded-full ${getStatusColor()}`} />
          <h3 className="text-xl font-semibold text-white">{container.name}</h3>
        </div>
        <div className="flex items-center gap-2">
          <div className={`w-2 h-2 rounded-full ${getConnectionStatusColor()}`} />
          <span className="text-xs text-gray-400">{status}</span>
        </div>
      </div>

      <div className="mb-4 space-y-2 text-sm text-gray-300">
        <div className="flex justify-between">
          <span>Image:</span>
          <span className="text-gray-400">{container.image}</span>
        </div>
        <div className="flex justify-between">
          <span>Status:</span>
          <span className="text-gray-400">{container.status}</span>
        </div>
        <div className="flex justify-between">
          <span>State:</span>
          <span className="text-gray-400 capitalize">{container.state}</span>
        </div>
      </div>

      {isVisible && <MetricsChart containerId={container.id} />}

      <div className="mt-4 flex gap-2">
        <button
          onClick={() => setShowDetail(true)}
          className="flex-1 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
        >
          View Details
        </button>
      </div>

      {showDetail && (
        <ContainerDetail
          containerId={container.id}
          containerName={container.name}
          onClose={() => setShowDetail(false)}
        />
      )}
    </div>
  );
}

