import { cn } from '@/lib/utils';

interface BadgePillProps {
  label: string;
  variant?: 'default' | 'orange' | 'pink' | 'violet' | 'emerald';
  className?: string;
}

const variants = {
  default: 'bg-surface-card text-ink',
  orange: 'bg-badge-orange text-ink',
  pink: 'bg-badge-pink text-ink',
  violet: 'bg-badge-violet text-ink',
  emerald: 'bg-badge-emerald text-ink',
};

export function BadgePill({ label, variant = 'default', className }: BadgePillProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center font-body text-caption font-medium',
        'rounded-pill px-3 py-1',
        variants[variant],
        className
      )}
    >
      {label}
    </span>
  );
}
