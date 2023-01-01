import React from 'react';
import ReactDOM from 'react-dom/client';

import './common/styles/index.css';
import { App, AppProvider } from './app';

const rootElement = document.getElementById('root') as HTMLElement;
ReactDOM.createRoot(rootElement).render(
  <React.StrictMode>
    <AppProvider>
      <App />
    </AppProvider>
  </React.StrictMode>
);
