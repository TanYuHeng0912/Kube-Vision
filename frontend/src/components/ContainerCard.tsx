import { useState, memo } from 'react';
import { useContainerStore, ContainerMetrics } from '../stores/containerStore';
import { useWebSocket } from '../hooks/useWebSocket';
import { useIntersectionObserver } from '../hooks/useIntersectionObserver';
import { WS_URL } from '../config';
import MetricsChart from './MetricsChart';
import ContainerDetail from './ContainerDetail';

interface ContainerCardProps {
  container: {
    id: string;
    name: string;
    image: string;
    status: string;
    state: string;
    labels?: Record<string, string>;
  };
}

function ContainerCard({ container }: ContainerCardProps) {
  const { addMetrics } = useContainerStore();
  const [cardRef, isVisible] = useIntersectionObserver<HTMLDivElement>();
  const [showDetail, setShowDetail] = useState(false);

  const wsUrl = isVisible
    ? `${WS_URL}/ws/stats/${container.id}`
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
        return 'bg-emerald-500';
      case 'exited':
        return 'bg-red-500';
      case 'paused':
        return 'bg-amber-500';
      default:
        return 'bg-gray-400';
    }
  };

  const getConnectionStatusColor = () => {
    switch (status) {
      case 'connected':
        return 'bg-emerald-500';
      case 'connecting':
        return 'bg-amber-400';
      case 'error':
      case 'disconnected':
        return 'bg-red-500';
      default:
        return 'bg-gray-400';
    }
  };

  const getStatusBadgeColor = () => {
    switch (container.state) {
      case 'running':
        return 'bg-emerald-50 text-emerald-700 border-emerald-200';
      case 'exited':
        return 'bg-red-50 text-red-700 border-red-200';
      case 'paused':
        return 'bg-amber-50 text-amber-700 border-amber-200';
      default:
        return 'bg-gray-50 text-gray-700 border-gray-200';
    }
  };

  return (
    <div
      ref={cardRef}
      className="bg-white rounded-xl p-6 shadow-md border border-gray-200 hover:shadow-lg transition-all duration-300"
    >
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-3">
          <div className={`w-3 h-3 rounded-full ${getStatusColor()}`} />
          <h3 className="text-xl font-semibold text-gray-800">{container.name}</h3>
        </div>
        <div className="flex items-center gap-2 px-3 py-1 rounded-full bg-blue-50 border border-blue-100">
          <div className={`w-2 h-2 rounded-full ${getConnectionStatusColor()}`} />
          <span className="text-xs text-blue-600 font-medium capitalize">{status}</span>
        </div>
      </div>

      <div className="mb-4 space-y-2 text-sm">
        <div className="flex justify-between py-2 px-3 rounded-lg bg-gray-50">
          <span className="text-gray-600">Image:</span>
          <span className="text-gray-800">{container.image}</span>
        </div>
        <div className="flex justify-between py-2 px-3 rounded-lg bg-gray-50">
          <span className="text-gray-600">Status:</span>
          <span className="text-gray-800">{container.status}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-gray-600">State:</span>
          <span className={`px-3 py-1 rounded-full text-xs font-bold border ${getStatusBadgeColor()}`}>
            {container.state}
          </span>
        </div>
      </div>

      {isVisible && <MetricsChart containerId={container.id} />}

      <div className="mt-4 flex gap-2">
        <button
          onClick={() => setShowDetail(true)}
          className="flex-1 px-4 py-2.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-all duration-200 font-semibold shadow-sm"
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

// Memoize component to prevent unnecessary re-renders
export default memo(ContainerCard, (prevProps, nextProps) => {
  // Only re-render if container properties change
  return (
    prevProps.container.id === nextProps.container.id &&
    prevProps.container.state === nextProps.container.state &&
    prevProps.container.status === nextProps.container.status
  );
});

