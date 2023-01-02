import React from 'react';
import { twMerge } from 'tailwind-merge';

export function Input({
  className,
  ...props
}: React.InputHTMLAttributes<HTMLInputElement>) {
  const baseClassLists = 'w-full rounded-md border py-3 px-5';
  const classLists = twMerge(baseClassLists, className);

  return <input {...props} className={classLists} />;
}
