import { useState } from 'react';
import LogTerminal from './LogTerminal';
import MetricsChart from './MetricsChart';

interface ContainerDetailProps {
  containerId: string;
  containerName: string;
  onClose: () => void;
}

export default function ContainerDetail({
  containerId,
  containerName,
  onClose,
}: ContainerDetailProps) {
  const [activeTab, setActiveTab] = useState<'metrics' | 'logs'>('metrics');

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg w-full max-w-6xl max-h-[90vh] flex flex-col shadow-lg">
        <div className="flex items-center justify-between p-4 border-b border-gray-200">
          <h2 className="text-2xl font-semibold text-gray-800">{containerName}</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 text-2xl"
          >
            ×
          </button>
        </div>

        <div className="flex border-b border-gray-200">
          <button
            onClick={() => setActiveTab('metrics')}
            className={`px-6 py-3 font-medium ${
              activeTab === 'metrics'
                ? 'text-blue-600 border-b-2 border-blue-600'
                : 'text-gray-500 hover:text-gray-700'
            }`}
          >
            Metrics
          </button>
          <button
            onClick={() => setActiveTab('logs')}
            className={`px-6 py-3 font-medium ${
              activeTab === 'logs'
                ? 'text-blue-600 border-b-2 border-blue-600'
                : 'text-gray-500 hover:text-gray-700'
            }`}
          >
            Logs
          </button>
        </div>

        <div className="flex-1 overflow-auto p-4">
          {activeTab === 'metrics' && (
            <div className="h-full">
              <MetricsChart containerId={containerId} showFullTime={true} />
            </div>
          )}
          {activeTab === 'logs' && (
            <div className="h-full">
              <LogTerminal containerId={containerId} visible={activeTab === 'logs'} />
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

