import { useEffect } from 'react';
import { useContainerStore } from './stores/containerStore';
import { api } from './services/api';
import ContainerCard from './components/ContainerCard';
import ErrorBoundary from './components/ErrorBoundary';

function App() {
  const { containers, setContainers, setLoading, setError } = useContainerStore();

  useEffect(() => {
    const fetchContainers = async () => {
      try {
        setLoading(true);
        setError(null);
        const data = await api.getContainers();
        setContainers(data);
      } catch (error) {
        setError(error instanceof Error ? error.message : 'Failed to fetch containers');
        console.error('Error fetching containers:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchContainers();
    
    // Refresh every 30 seconds
    const interval = setInterval(fetchContainers, 30000);
    return () => clearInterval(interval);
  }, [setContainers, setLoading, setError]);

  const { isLoading, error } = useContainerStore();

  if (isLoading && containers.length === 0) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-blue-100 flex items-center justify-center">
        <div className="text-blue-700 text-xl font-semibold">Loading containers...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-blue-100 flex items-center justify-center">
        <div className="text-red-600 text-xl font-semibold">Error: {error}</div>
      </div>
    );
  }

  return (
    <ErrorBoundary>
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-blue-50 p-8">
      <div className="max-w-7xl mx-auto">
        <header className="mb-8">
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-gradient-to-br from-blue-600 to-cyan-500 rounded-lg flex items-center justify-center shadow-md">
                <div className="flex items-end gap-1 h-6">
                  <div className="w-1.5 bg-white rounded-sm" style={{ height: '40%' }}></div>
                  <div className="w-1.5 bg-white rounded-sm" style={{ height: '70%' }}></div>
                  <div className="w-1.5 bg-white rounded-sm" style={{ height: '100%' }}></div>
                </div>
              </div>
              <div>
                <h1 className="text-5xl font-extrabold bg-gradient-to-r from-blue-700 via-blue-600 to-cyan-500 bg-clip-text text-transparent mb-2 drop-shadow-sm">KubeVision</h1>
                <p className="text-gray-600 font-medium">Docker Container Monitoring Dashboard</p>
              </div>
            </div>
          </div>
        </header>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {containers.map((container) => (
            <ContainerCard key={container.id} container={container} />
          ))}
        </div>

        {containers.length === 0 && (
          <div className="text-center text-gray-500 mt-12 font-medium">
            No containers found. Start some Docker containers to see them here.
          </div>
        )}
      </div>
    </div>
    </ErrorBoundary>
  );
}

export default App;

