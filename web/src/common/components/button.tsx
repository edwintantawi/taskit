import React from 'react';
import { twMerge } from 'tailwind-merge';
import { To as LinkPropsTo } from 'react-router-dom';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  fullWidth?: boolean;
  size?: 'small' | 'medium' | 'large';
  variants?: 'normal' | 'contained' | 'outlined';
  as?: React.ElementType;
  to?: LinkPropsTo;
}

export function Button({
  className,
  fullWidth = false,
  size = 'medium',
  variants = 'normal',
  as = 'button',
  ...props
}: ButtonProps) {
  const sizeStyle = {
    small: 'py-2 px-3 text-xs',
    medium: 'py-2 px-4 text-sm',
    large: 'py-3 px-5 text-sm md:text-base',
  };

  const variantStyle = {
    normal: 'text-gray-900 hover:bg-gray-100',
    contained:
      'bg-gray-900 text-white hover:bg-gray-800 border border-gray-900',
    outlined: 'bg-white text-gray-900 border border-gray-900 hover:bg-gray-100',
  };

  const Component = as;

  return (
    <Component
      {...props}
      className={twMerge(
        'inline-block rounded-md disabled:cursor-not-allowed',
        fullWidth && 'w-full',
        sizeStyle[size],
        variantStyle[variants],
        className
      )}
    />
  );
}
