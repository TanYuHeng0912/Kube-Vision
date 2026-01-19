import { useEffect } from 'react';
import { useContainerStore } from './stores/containerStore';
import { api } from './services/api';
import ContainerCard from './components/ContainerCard';

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
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <div className="text-white text-xl">Loading containers...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <div className="text-red-500 text-xl">Error: {error}</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-900 p-8">
      <div className="max-w-7xl mx-auto">
        <header className="mb-8">
          <h1 className="text-4xl font-bold text-white mb-2">KubeVision</h1>
          <p className="text-gray-400">Docker Container Monitoring Dashboard</p>
        </header>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {containers.map((container) => (
            <ContainerCard key={container.id} container={container} />
          ))}
        </div>

        {containers.length === 0 && (
          <div className="text-center text-gray-400 mt-12">
            No containers found. Start some Docker containers to see them here.
          </div>
        )}
      </div>
    </div>
  );
}

export default App;

