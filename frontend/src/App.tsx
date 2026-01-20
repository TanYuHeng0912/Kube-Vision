import { useEffect, useState, useMemo } from 'react';
import { useContainerStore } from './stores/containerStore';
import { api, ContainerInfo } from './services/api';
import ContainerCard from './components/ContainerCard';
import ErrorBoundary from './components/ErrorBoundary';
import SearchAndFilter from './components/SearchAndFilter';
import ImageList from './components/ImageList';
import EventsStream from './components/EventsStream';
import { groupContainers } from './utils/grouping';

function App() {
  const { containers, setContainers, setLoading, setError } = useContainerStore();
  const [filteredContainers, setFilteredContainers] = useState<ContainerInfo[]>([]);
  const [groupBy, setGroupBy] = useState<string>('none');
  const [activeTab, setActiveTab] = useState<'containers' | 'images' | 'events'>('containers');

  useEffect(() => {
    const fetchContainers = async () => {
      try {
        setLoading(true);
        setError(null);
        const data = await api.getContainers();
        setContainers(data);
        setFilteredContainers(data);
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

  const groupedContainers = useMemo(() => {
    return groupContainers(filteredContainers, groupBy);
  }, [filteredContainers, groupBy]);

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

  const renderContainers = () => {
    if (groupBy === 'none') {
      return (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredContainers.map((container) => (
            <ContainerCard key={container.id} container={container} />
          ))}
        </div>
      );
    }

    return (
      <div className="space-y-8">
        {Object.entries(groupedContainers).map(([groupName, groupContainers]) => (
          <div key={groupName}>
            <h2 className="text-2xl font-bold text-gray-800 mb-4 pb-2 border-b border-gray-200">
              {groupName} <span className="text-sm font-normal text-gray-500">({groupContainers.length})</span>
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {groupContainers.map((container) => (
                <ContainerCard key={container.id} container={container} />
              ))}
            </div>
          </div>
        ))}
      </div>
    );
  };

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

        {/* Tabs */}
        <div className="mb-6 border-b border-gray-200">
          <nav className="flex gap-4">
            <button
              onClick={() => setActiveTab('containers')}
              className={`px-4 py-2 font-medium border-b-2 transition-colors ${
                activeTab === 'containers'
                  ? 'border-blue-600 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              Containers
            </button>
            <button
              onClick={() => setActiveTab('images')}
              className={`px-4 py-2 font-medium border-b-2 transition-colors ${
                activeTab === 'images'
                  ? 'border-blue-600 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              Images
            </button>
            <button
              onClick={() => setActiveTab('events')}
              className={`px-4 py-2 font-medium border-b-2 transition-colors ${
                activeTab === 'events'
                  ? 'border-blue-600 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              Events
            </button>
          </nav>
        </div>

        {/* Containers Tab */}
        {activeTab === 'containers' && (
          <>
            <SearchAndFilter
              containers={containers}
              onFiltered={setFilteredContainers}
              onGroupChange={setGroupBy}
            />
            {renderContainers()}
            {filteredContainers.length === 0 && containers.length > 0 && (
              <div className="text-center text-gray-500 mt-12 font-medium">
                No containers match your filters.
              </div>
            )}
            {containers.length === 0 && (
              <div className="text-center text-gray-500 mt-12 font-medium">
                No containers found. Start some Docker containers to see them here.
              </div>
            )}
          </>
        )}

        {/* Images Tab */}
        {activeTab === 'images' && <ImageList />}

        {/* Events Tab */}
        {activeTab === 'events' && <EventsStream />}
      </div>
    </div>
    </ErrorBoundary>
  );
}

export default App;

