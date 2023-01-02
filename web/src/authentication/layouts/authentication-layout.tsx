import React from 'react';

interface AuthenticationLayoutProps {
  children?: React.ReactNode;
  form: React.ReactNode;
  title: string;
  subtitle: string;
}

export function AuthenticationLayout({
  children,
  form,
  title,
  subtitle,
}: AuthenticationLayoutProps) {
  return (
    <div className="rounded-md border p-8">
      <header className="mb-8 space-y-1">
        <h1 className="text-center text-2xl font-bold uppercase">{title}</h1>
        <p className="text-center text-gray-500">{subtitle}</p>
      </header>
      {form}
      <div className="pt-6">{children}</div>
    </div>
  );
}
