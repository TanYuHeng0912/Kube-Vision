# KubeVision Frontend

React + TypeScript frontend for KubeVision Docker monitoring dashboard.

## Features

- Real-time container metrics visualization with ECharts
- WebSocket-based streaming
- Virtual scrolling with Intersection Observer
- Zustand state management
- TailwindCSS styling

## Prerequisites

- Node.js 18+ and npm

## Installation

```bash
npm install
```

## Development

```bash
npm run dev
```

The app will be available at `http://localhost:5173`

## Build

```bash
npm run build
```

## Configuration

Create a `.env` file:

```env
VITE_API_URL=http://localhost:8080
```

