import React from 'react';
import { twMerge } from 'tailwind-merge';

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  error?: string;
}

export const Input = React.forwardRef(function Input(
  { className, error, ...props }: InputProps,
  ref: React.ForwardedRef<HTMLInputElement>
) {
  const baseClassLists =
    'w-full rounded-md border py-3 px-5 focus:outline-gray-900 placeholder:text-sm md:placeholder:text-base';
  const classLists = twMerge(baseClassLists, className);

  return (
    <div>
      <input {...props} className={classLists} ref={ref} />
      {error && <span className="m-2 block text-sm text-red-500">{error}</span>}
    </div>
  );
});
