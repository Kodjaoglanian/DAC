'use client';

import { cn } from '@/lib/utils';
import { ButtonHTMLAttributes, forwardRef } from 'react';

interface ButtonPrimaryProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  isLoading?: boolean;
}

export const ButtonPrimary = forwardRef<HTMLButtonElement, ButtonPrimaryProps>(
  ({ className, children, isLoading, disabled, ...props }, ref) => {
    return (
      <button
        ref={ref}
        className={cn(
          'inline-flex items-center justify-center h-10 px-5 py-3',
          'bg-primary text-on-primary font-body text-sm font-semibold leading-none',
          'rounded-md transition-colors',
          'active:bg-primary-active',
          'disabled:bg-primary-disabled disabled:text-muted disabled:cursor-not-allowed',
          className
        )}
        disabled={disabled || isLoading}
        {...props}
      >
        {isLoading ? (
          <div className="w-4 h-4 border-2 border-on-primary border-t-transparent rounded-full animate-spin" />
        ) : (
          children
        )}
      </button>
    );
  }
);
ButtonPrimary.displayName = 'ButtonPrimary';
