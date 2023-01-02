import React from 'react';

interface MainLayoutProps {
  children?: React.ReactNode;
}

export function MainLayout({ children }: MainLayoutProps) {
  return <main className="flex flex-1 flex-col pt-6">{children}</main>;
}
