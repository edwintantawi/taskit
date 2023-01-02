import React from 'react';

interface MainLayoutProps {
  children?: React.ReactNode;
}

export function MainLayout({ children }: MainLayoutProps) {
  return <main className="flex-1 py-6">{children}</main>;
}
