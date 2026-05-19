'use client';

import { cn } from '@/lib/utils';
import { InputHTMLAttributes, forwardRef } from 'react';

export const TextInput = forwardRef<HTMLInputElement, InputHTMLAttributes<HTMLInputElement>>(
  ({ className, ...props }, ref) => {
    return (
      <input
        ref={ref}
        className={cn(
          'w-full h-10 bg-canvas text-ink font-body text-base',
          'rounded-md px-3.5 py-2.5',
          'border border-hairline',
          'focus:border-ink focus:outline-none',
          'placeholder:text-muted',
          'transition-colors',
          className
        )}
        {...props}
      />
    );
  }
);
TextInput.displayName = 'TextInput';
