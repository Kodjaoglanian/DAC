'use client';

import { Toaster } from 'sonner';

export function ToastProvider() {
  return (
    <Toaster
      position="top-right"
      toastOptions={{
        style: {
          background: '#111111',
          color: '#ffffff',
          border: 'none',
        },
      }}
    />
  );
}
