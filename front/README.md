# Enter Leave System - Frontend

This is the frontend application for the Enter Leave System, built with React 19, TypeScript, and Vite.

## Technology Stack

- **React 19** - UI library
- **TypeScript** - Type-safe development
- **Vite** - Fast build tool and development server
- **React Router v7** - Client-side routing
- **React Bootstrap** - UI components
- **CSS Modules** - Scoped styling

## Prerequisites

- Node.js 20.x or later
- npm (comes with Node.js)

## Development Setup

### Install Dependencies

```bash
npm install
```

### Environment Variables

Create or modify the `.env` file in the root of the `front/` directory:

```env
VITE_ROOMNAME=Aizu Geek Dojo
VITE_SHOWQUESTION=true
```

Available environment variables:
- `VITE_ROOMNAME` - The name of the room displayed on the top page
- `VITE_SHOWQUESTION` - Whether to show the question page on exit (`true` or `false`)

### Start Development Server

```bash
npm run dev
```

The application will be available at `http://localhost:5173` by default.

## Building for Production

### Build for Development/Testing

```bash
npm run build
```

The built files will be output to the `dist/` directory.

## Project Structure

```
front/
├── public/              # Static assets
├── src/
│   ├── components/      # React components
│   │   ├── Top.tsx              # Main page with card reader
│   │   ├── Register.tsx         # Card registration page
│   │   ├── Welcome.tsx          # Welcome/check-in page
│   │   ├── Goodbye.tsx          # Goodbye/check-out page
│   │   ├── Question.tsx         # Purpose questionnaire page
│   │   ├── Forgot.tsx           # Forgot card page
│   │   └── *.module.css         # Component styles (CSS Modules)
│   ├── types/           # TypeScript type definitions
│   │   └── index.ts
│   ├── utils/           # Utility functions
│   │   └── api.ts               # API communication utilities
│   ├── App.tsx          # Root component with router
│   ├── App.module.css   # App component styles
│   ├── main.tsx         # Application entry point
│   └── index.css        # Global styles
├── .env                 # Development environment variables
├── .env.agd             # AGD production environment variables
├── .env.gl              # GL production environment variables
├── index.html           # HTML template
├── package.json         # Dependencies and scripts
├── tsconfig.json        # TypeScript configuration
├── tsconfig.node.json   # TypeScript configuration for Node
└── vite.config.ts       # Vite configuration
```

## Key Components

### Top (Top.tsx)
- Main page that displays instructions to scan card
- Establishes WebSocket connection to card reader
- Handles card reading events and routes to appropriate pages
- Includes "Forgot card?" link

### Register (Register.tsx)
- Card registration page for new cards
- Accepts student ID input
- Validates and registers card with backend API

### Welcome (Welcome.tsx)
- Welcome page displayed on check-in
- Shows personalized greeting
- Automatically redirects to top page after 5 seconds
- Records check-in log via API

### Goodbye (Goodbye.tsx)
- Goodbye page displayed on check-out
- Simple farewell message
- Automatically redirects to top page after 3 seconds

### Question (Question.tsx)
- Questionnaire page for exit purposes
- Collects usage purpose (3D Printer, Laser Cutter, Training, Other)
- Optional text field for additional requests
- Supports Ctrl+Enter submission
- Records response via API
- Can be disabled via `VITE_SHOWQUESTION` environment variable

### Forgot (Forgot.tsx)
- Manual check-in/out for users who forgot their card
- Accepts student ID input
- Routes to appropriate page based on current status

## API Integration

The frontend communicates with a backend server running on `localhost:3000`:

- **WebSocket**: `ws://localhost:3000/socket/readCard` - Real-time card reader events
- **GET** `/api/user?sid={sid}` - Fetch user information
- **POST** `/api/user` - Register card
- **POST** `/api/log` - Record check-in/check-out logs

## Deployment

The application is automatically built and deployed via GitHub Actions CI/CD pipeline. See [.github/workflows/main.yml](../.github/workflows/main.yml) for details.

## Notes

- The application requires a running backend server for full functionality
- WebSocket connection to card reader is essential for automatic card scanning
- CSS Modules are used for component-scoped styling
- All timestamps are handled by the backend server
