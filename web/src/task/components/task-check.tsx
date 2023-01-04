import React from 'react';
import { CheckIcon } from '@heroicons/react/24/outline';

interface TaskCheckProps {
  onChange?: (event: React.InputHTMLAttributes<HTMLInputElement>) => void;
  checked?: boolean;
}

export function TaskCheck({ checked, onChange }: TaskCheckProps) {
  return (
    <div className="relative mt-1">
      <input
        type="checkbox"
        className="peer absolute top-0 z-50 aspect-square h-5 w-5 cursor-pointer opacity-0"
        checked={checked}
        onChange={onChange}
      />
      <div
        className={`grid h-5 w-5 place-items-center rounded-full border border-gray-400 text-white ${
          checked ? 'bg-gray-900' : 'bg-white peer-hover:text-gray-400'
        }`}
      >
        <CheckIcon className="h-3 w-3" />
      </div>
    </div>
  );
}
