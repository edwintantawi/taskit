import React from 'react';
import { BrowserRouter } from 'react-router-dom';

interface AppProviderProps {
  children?: React.ReactNode;
}

// Setup all application provider
export function AppProvider({ children }: AppProviderProps) {
  return <BrowserRouter>{children}</BrowserRouter>;
}
