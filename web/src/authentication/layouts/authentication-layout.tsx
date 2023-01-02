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
    <div className="rounded-md border p-5 md:p-8">
      <header className="mb-8 space-y-1">
        <h1 className="text-center text-xl font-bold uppercase md:text-2xl">
          {title}
        </h1>
        <p className="text-center text-sm text-gray-500 md:text-base">
          {subtitle}
        </p>
      </header>
      {form}
      <div className="pt-6">{children}</div>
    </div>
  );
}
