'use client';

import { useEffect } from 'react';

export default function Error({ error, reset }: { error: Error & { digest?: string }; reset: () => void }) {
  useEffect(() => {
    console.error('Next.js Error Boundary:', error);
  }, [error]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-canvas">
      <div className="text-center">
        <h2 className="text-display-sm text-ink mb-4">Erro no carregamento</h2>
        <p className="text-muted mb-4">{error.message}</p>
        <button
          onClick={reset}
          className="bg-primary text-on-primary px-4 py-2 rounded-md font-body text-sm font-semibold"
        >
          Tentar novamente
        </button>
      </div>
    </div>
  );
}
