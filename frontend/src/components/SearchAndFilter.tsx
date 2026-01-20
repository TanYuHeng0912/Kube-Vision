import { useState, useMemo } from 'react';
import { ContainerInfo } from '../services/api';

interface SearchAndFilterProps {
  containers: ContainerInfo[];
  onFiltered: (filtered: ContainerInfo[]) => void;
  onGroupChange: (groupBy: string) => void;
}

export default function SearchAndFilter({ containers, onFiltered, onGroupChange }: SearchAndFilterProps) {
  const [searchTerm, setSearchTerm] = useState('');
  const [stateFilter, setStateFilter] = useState<string>('all');
  const [groupBy, setGroupBy] = useState<string>('none');

  const filteredContainers = useMemo(() => {
    let filtered = [...containers];

    // Search filter
    if (searchTerm) {
      const term = searchTerm.toLowerCase();
      filtered = filtered.filter(
        (c) =>
          c.name.toLowerCase().includes(term) ||
          c.image.toLowerCase().includes(term) ||
          c.id.toLowerCase().includes(term) ||
          c.status.toLowerCase().includes(term)
      );
    }

    // State filter
    if (stateFilter !== 'all') {
      filtered = filtered.filter((c) => c.state === stateFilter);
    }

    onFiltered(filtered);
    return filtered;
  }, [containers, searchTerm, stateFilter, onFiltered]);

  const handleGroupChange = (value: string) => {
    setGroupBy(value);
    onGroupChange(value);
  };

  return (
    <div className="mb-6 space-y-4">
      <div className="flex flex-wrap gap-4 items-center">
        {/* Search */}
        <div className="flex-1 min-w-[200px]">
          <div className="relative">
            <input
              type="text"
              placeholder="Search containers by name, image, or ID..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full px-4 py-2 pl-10 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
            <svg
              className="absolute left-3 top-2.5 w-5 h-5 text-gray-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
              />
            </svg>
          </div>
        </div>

        {/* State Filter */}
        <div className="flex items-center gap-2">
          <label className="text-sm font-medium text-gray-700">State:</label>
          <select
            value={stateFilter}
            onChange={(e) => setStateFilter(e.target.value)}
            className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="all">All</option>
            <option value="running">Running</option>
            <option value="exited">Exited</option>
            <option value="paused">Paused</option>
            <option value="created">Created</option>
          </select>
        </div>

        {/* Group By */}
        <div className="flex items-center gap-2">
          <label className="text-sm font-medium text-gray-700">Group:</label>
          <select
            value={groupBy}
            onChange={(e) => handleGroupChange(e.target.value)}
            className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="none">None</option>
            <option value="project">By Project</option>
            <option value="environment">By Environment</option>
            <option value="image">By Image</option>
            <option value="state">By State</option>
          </select>
        </div>
      </div>

      {/* Results count */}
      <div className="text-sm text-gray-600">
        Showing {filteredContainers.length} of {containers.length} containers
      </div>
    </div>
  );
}





