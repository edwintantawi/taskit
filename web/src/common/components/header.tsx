import React from 'react';

interface HeaderProps {
  title?: string;
  subtitle?: string;
}

export function Header({ title, subtitle }: HeaderProps) {
  return (
    <header className="text-center">
      <h1 className="mx-auto mb-3 max-w-[500px] text-2xl font-bold md:text-3xl">
        {title}
      </h1>
      {subtitle && (
        <p className="mx-auto mb-8 max-w-[300px] text-sm text-gray-500 md:text-base">
          {subtitle}
        </p>
      )}
    </header>
  );
}
