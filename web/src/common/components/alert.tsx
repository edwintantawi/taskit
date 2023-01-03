import React from 'react';
import { twMerge } from 'tailwind-merge';
import {
  ExclamationCircleIcon,
  InformationCircleIcon,
} from '@heroicons/react/24/outline';

interface AlertProps {
  children?: React.ReactNode;
  severity?: 'info' | 'error';
}

export function Alert({ children, severity = 'info' }: AlertProps) {
  const severityStyle = {
    info: 'bg-blue-50 border border-blue-600 text-blue-600',
    error: 'bg-red-50 border border-red-600 text-red-600',
  };

  const severityIcon = {
    info: <InformationCircleIcon className="h-6 w-6" />,
    error: <ExclamationCircleIcon className="h-6 w-6" />,
  };

  return (
    <div
      className={twMerge(
        'flex items-start gap-3 rounded-md py-3 px-4 text-sm leading-6',
        severityStyle[severity]
      )}
    >
      <div>{severityIcon[severity]}</div>
      {children}
    </div>
  );
}
