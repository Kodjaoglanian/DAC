'use client';

import { cn } from '@/lib/utils';
import { ButtonHTMLAttributes, forwardRef } from 'react';

export const ButtonSecondary = forwardRef<HTMLButtonElement, ButtonHTMLAttributes<HTMLButtonElement>>(
  ({ className, children, ...props }, ref) => {
    return (
      <button
        ref={ref}
        className={cn(
          'inline-flex items-center justify-center h-10 px-5 py-3',
          'bg-canvas text-ink font-body text-sm font-semibold leading-none',
          'rounded-md border border-hairline transition-colors',
          'active:bg-surface-soft',
          className
        )}
        {...props}
      >
        {children}
      </button>
    );
  }
);
ButtonSecondary.displayName = 'ButtonSecondary';
