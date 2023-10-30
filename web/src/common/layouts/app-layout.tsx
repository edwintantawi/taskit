import React from 'react';

interface AppLayoutProps {
  children?: React.ReactNode;
}

export function AppLayout({ children }: AppLayoutProps) {
  return (
    <div className="mx-auto flex min-h-screen max-w-2xl flex-col px-4">
      {children}
    </div>
  );
}
