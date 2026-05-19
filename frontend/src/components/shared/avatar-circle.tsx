import { cn } from '@/lib/utils';
import { getInitials } from '@/lib/utils';

interface AvatarCircleProps {
  name: string;
  src?: string | null;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

const sizes = {
  sm: 'w-8 h-8 text-xs',
  md: 'w-9 h-9 text-caption',
  lg: 'w-12 h-12 text-sm',
};

export function AvatarCircle({ name, src, size = 'md', className }: AvatarCircleProps) {
  return (
    <div
      className={cn(
        'rounded-full bg-surface-card flex items-center justify-center overflow-hidden',
        sizes[size],
        className
      )}
    >
      {src ? (
        <img src={src} alt={name} className="w-full h-full object-cover" />
      ) : (
        <span className="font-body font-medium text-ink">{getInitials(name)}</span>
      )}
    </div>
  );
}
