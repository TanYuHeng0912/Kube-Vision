import { useEffect, useState, useRef } from 'react';
import { useWebSocket } from '../hooks/useWebSocket';
import { WS_URL } from '../config';

interface DockerEvent {
  type: string;
  action: string;
  actor: {
    ID?: string;
    Attributes?: Record<string, string>;
  };
  time: number;
  timeNano: number;
}

export default function EventsStream() {
  const [events, setEvents] = useState<DockerEvent[]>([]);
  const [filter, setFilter] = useState<string>('all');
  const [autoScroll, setAutoScroll] = useState(true);
  const eventsEndRef = useRef<HTMLDivElement>(null);

  const wsUrl = `${WS_URL}/ws/events`;

  const { status } = useWebSocket({
    url: wsUrl,
    onMessage: (data: DockerEvent) => {
      setEvents((prev) => {
        const newEvents = [data, ...prev];
        return newEvents.slice(0, 500); // Keep last 500 events
      });
    },
    reconnect: true,
  });

  useEffect(() => {
    if (autoScroll && eventsEndRef.current) {
      eventsEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [events, autoScroll]);

  const filteredEvents = events.filter((event) => {
    if (filter === 'all') return true;
    return event.type === filter;
  });

  const formatTime = (timestamp: number): string => {
    const date = new Date(timestamp * 1000);
    return date.toLocaleString();
  };

  const getEventColor = (type: string, action: string): string => {
    if (type === 'container') {
      if (action.includes('start')) return 'text-emerald-600';
      if (action.includes('stop') || action.includes('die')) return 'text-red-600';
      if (action.includes('create')) return 'text-blue-600';
      return 'text-gray-600';
    }
    if (type === 'image') {
      return 'text-purple-600';
    }
    return 'text-gray-600';
  };

  const getEventIcon = (type: string, action: string): string => {
    if (type === 'container') {
      if (action.includes('start')) return '▶';
      if (action.includes('stop') || action.includes('die')) return '⏹';
      if (action.includes('create')) return '➕';
      return '📦';
    }
    if (type === 'image') {
      return '🖼️';
    }
    return '📋';
  };

  return (
    <div className="space-y-4">
      {/* Controls */}
      <div className="flex items-center justify-between bg-white rounded-lg p-4 shadow-md">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <div className={`w-3 h-3 rounded-full ${status === 'connected' ? 'bg-emerald-500' : 'bg-red-500'}`} />
            <span className="text-sm text-gray-600 capitalize">{status}</span>
          </div>
          <div className="flex items-center gap-2">
            <label className="text-sm font-medium text-gray-700">Filter:</label>
            <select
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              className="px-3 py-1 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="all">All Events</option>
              <option value="container">Container</option>
              <option value="image">Image</option>
              <option value="network">Network</option>
              <option value="volume">Volume</option>
            </select>
          </div>
        </div>
        <div className="flex items-center gap-4">
          <label className="flex items-center gap-2 text-sm text-gray-700">
            <input
              type="checkbox"
              checked={autoScroll}
              onChange={(e) => setAutoScroll(e.target.checked)}
              className="rounded"
            />
            Auto-scroll
          </label>
          <button
            onClick={() => setEvents([])}
            className="px-3 py-1 text-sm bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200"
          >
            Clear
          </button>
        </div>
      </div>

      {/* Events List */}
      <div className="bg-white rounded-xl shadow-md p-4 max-h-[600px] overflow-y-auto">
        {filteredEvents.length === 0 ? (
          <div className="text-center text-gray-500 py-8">
            {status === 'connected' ? 'No events yet. Waiting for Docker events...' : 'Connecting to events stream...'}
          </div>
        ) : (
          <div className="space-y-2">
            {filteredEvents.map((event, index) => (
              <div
                key={`${event.timeNano}-${index}`}
                className="flex items-start gap-3 p-3 rounded-lg hover:bg-gray-50 border-l-4 border-gray-200"
              >
                <div className="text-2xl">{getEventIcon(event.type, event.action)}</div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <span className={`font-semibold ${getEventColor(event.type, event.action)}`}>
                      {event.type}
                    </span>
                    <span className="text-gray-600">•</span>
                    <span className="text-gray-700">{event.action}</span>
                  </div>
                  {event.actor?.ID && (
                    <div className="text-sm text-gray-500 mt-1">
                      ID: <span className="font-mono">{event.actor.ID.substring(0, 12)}</span>
                    </div>
                  )}
                  {event.actor?.Attributes && Object.keys(event.actor.Attributes).length > 0 && (
                    <div className="text-xs text-gray-400 mt-1">
                      {Object.entries(event.actor.Attributes)
                        .slice(0, 3)
                        .map(([key, value]) => `${key}=${value}`)
                        .join(', ')}
                    </div>
                  )}
                  <div className="text-xs text-gray-400 mt-1">{formatTime(event.time)}</div>
                </div>
              </div>
            ))}
            <div ref={eventsEndRef} />
          </div>
        )}
      </div>

      <div className="text-sm text-gray-600">
        Showing {filteredEvents.length} of {events.length} events
      </div>
    </div>
  );
}


