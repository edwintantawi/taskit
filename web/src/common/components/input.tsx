import React from 'react';
import { twMerge } from 'tailwind-merge';

export function Input({
  className,
  ...props
}: React.InputHTMLAttributes<HTMLInputElement>) {
  const baseClassLists =
    'w-full rounded-md border py-3 px-5 focus:outline-gray-900 placeholder:text-sm md:placeholder:text-base';
  const classLists = twMerge(baseClassLists, className);

  return <input {...props} className={classLists} />;
}
