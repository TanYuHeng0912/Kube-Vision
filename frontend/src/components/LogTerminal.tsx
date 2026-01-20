import { useEffect, useRef } from 'react';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { SearchAddon } from 'xterm-addon-search';
import 'xterm/css/xterm.css';
import { useWebSocket } from '../hooks/useWebSocket';
import { WS_URL } from '../config';

interface LogTerminalProps {
  containerId: string;
  visible?: boolean;
}

export default function LogTerminal({ containerId, visible = true }: LogTerminalProps) {
  const terminalRef = useRef<HTMLDivElement>(null);
  const terminalInstanceRef = useRef<Terminal | null>(null);
  const fitAddonRef = useRef<FitAddon | null>(null);
  const searchAddonRef = useRef<SearchAddon | null>(null);
  const autoScrollRef = useRef(true); // Use ref to track auto-scroll state without re-renders

  const wsUrl = visible
    ? `${WS_URL}/ws/logs/${containerId}?follow=true&tail=100`
    : null;

  const { status } = useWebSocket({
    url: wsUrl,
    onMessage: (data: any) => {
      if (terminalInstanceRef.current) {
        // Handle both string and ArrayBuffer
        if (typeof data === 'string') {
          terminalInstanceRef.current.write(data);
        } else if (data instanceof ArrayBuffer) {
          const decoder = new TextDecoder();
          terminalInstanceRef.current.write(decoder.decode(data));
        }
      }
    },
    onError: (error) => {
      console.error('WebSocket error:', error);
      if (terminalInstanceRef.current) {
        terminalInstanceRef.current.writeln('\r\n\x1b[31m[Connection Error]\x1b[0m');
      }
    },
    reconnect: true,
  });

  useEffect(() => {
    if (!terminalRef.current || !visible) return;

    // Initialize terminal
    const terminal = new Terminal({
      cursorBlink: true,
      fontSize: 14,
      fontFamily: 'Consolas, "Courier New", monospace',
      theme: {
        background: '#ffffff',
        foreground: '#333333',
        cursor: '#333333',
        black: '#000000',
        red: '#cd3131',
        green: '#0dbc79',
        yellow: '#e5e510',
        blue: '#2472c8',
        magenta: '#bc3fbc',
        cyan: '#11a8cd',
        white: '#e5e5e5',
        brightBlack: '#666666',
        brightRed: '#f14c4c',
        brightGreen: '#23d18b',
        brightYellow: '#f5f543',
        brightBlue: '#3b8eea',
        brightMagenta: '#d670d6',
        brightCyan: '#29b8db',
        brightWhite: '#e5e5e5',
      },
    });

    // Add addons
    const fitAddon = new FitAddon();
    const searchAddon = new SearchAddon();
    
    terminal.loadAddon(fitAddon);
    terminal.loadAddon(searchAddon);

    // Open terminal
    terminal.open(terminalRef.current);
    fitAddon.fit();

    // Store references
    terminalInstanceRef.current = terminal;
    fitAddonRef.current = fitAddon;
    searchAddonRef.current = searchAddon;

    // Handle resize
    const handleResize = () => {
      if (fitAddonRef.current) {
        fitAddonRef.current.fit();
      }
    };
    window.addEventListener('resize', handleResize);

    // Handle copy/paste
    terminal.attachCustomKeyEventHandler((event) => {
      // Ctrl+C
      if (event.ctrlKey && event.key === 'c' && terminal.hasSelection()) {
        const selection = terminal.getSelection();
        if (selection) {
          navigator.clipboard.writeText(selection);
        }
        return false;
      }
      // Ctrl+V
      if (event.ctrlKey && event.key === 'v') {
        navigator.clipboard.readText().then((text) => {
          terminal.paste(text);
        });
        return false;
      }
      return true;
    });

    // Auto-scroll handling - use ref to avoid re-creating terminal
    const handleScroll = () => {
      if (terminalRef.current) {
        const element = terminalRef.current.querySelector('.xterm-viewport');
        if (element) {
          const isAtBottom =
            element.scrollTop + element.clientHeight >= element.scrollHeight - 10;
          autoScrollRef.current = isAtBottom;
        }
      }
    };

    if (terminalRef.current) {
      const viewport = terminalRef.current.querySelector('.xterm-viewport');
      if (viewport) {
        viewport.addEventListener('scroll', handleScroll);
      }
    }

    // Auto-scroll to bottom when new data arrives
    const originalWrite = terminal.write.bind(terminal);
    terminal.write = function (data: string) {
      originalWrite(data);
      if (autoScrollRef.current && terminalRef.current) {
        const viewport = terminalRef.current.querySelector('.xterm-viewport');
        if (viewport) {
          viewport.scrollTop = viewport.scrollHeight;
        }
      }
    };

    return () => {
      window.removeEventListener('resize', handleResize);
      if (terminalRef.current) {
        const viewport = terminalRef.current.querySelector('.xterm-viewport');
        if (viewport) {
          viewport.removeEventListener('scroll', handleScroll);
        }
      }
      terminal.dispose();
      terminalInstanceRef.current = null;
    };
  }, [visible]); // Removed autoScroll from dependencies - this was causing the issue!

  const handleClear = () => {
    if (terminalInstanceRef.current) {
      terminalInstanceRef.current.clear();
    }
  };

  const handleSearch = (query: string) => {
    if (searchAddonRef.current && query) {
      searchAddonRef.current.findNext(query);
    }
  };

  return (
    <div className="flex flex-col h-full bg-white rounded-lg overflow-hidden border border-gray-200 shadow-sm">
      <div className="flex items-center justify-between p-2 bg-gray-50 border-b border-gray-200">
        <div className="flex items-center gap-2">
          <div
            className={`w-2 h-2 rounded-full ${
              status === 'connected'
                ? 'bg-green-500'
                : status === 'connecting'
                ? 'bg-yellow-500'
                : 'bg-red-500'
            }`}
          />
          <span className="text-xs text-gray-600">Logs - {status}</span>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={handleClear}
            className="px-2 py-1 text-xs bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Clear
          </button>
          <input
            type="text"
            placeholder="Search logs..."
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleSearch(e.currentTarget.value);
              }
            }}
            className="px-2 py-1 text-xs bg-white text-gray-800 rounded border border-gray-300 focus:outline-none focus:border-blue-500"
          />
        </div>
      </div>
      <div ref={terminalRef} className="flex-1 p-2" style={{ minHeight: '400px' }} />
    </div>
  );
}

