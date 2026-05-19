'use client';

import Link from 'next/link';
import { useAuth } from '@/providers/auth-provider';
import { AvatarCircle } from '@/components/shared/avatar-circle';
import { LayoutDashboard, FolderKanban, BarChart3, Settings, LogOut, Menu } from 'lucide-react';
import { useState } from 'react';

export function TopNav() {
  const { user, logout } = useAuth();
  const [showMenu, setShowMenu] = useState(false);

  return (
    <header className="h-16 bg-canvas border-b border-hairline-soft flex items-center justify-between px-4 lg:px-8 sticky top-0 z-50">
      <div className="flex items-center gap-8">
        <Link href="/dashboard" className="font-display text-title-lg text-ink">
          ProjectTracker
        </Link>

        <nav className="hidden md:flex items-center gap-1">
          <NavLink href="/dashboard" icon={<LayoutDashboard className="w-4 h-4" />}>Dashboard</NavLink>
          <NavLink href="/projects" icon={<FolderKanban className="w-4 h-4" />}>Projetos</NavLink>
          <NavLink href="/reports" icon={<BarChart3 className="w-4 h-4" />}>Relatórios</NavLink>
        </nav>
      </div>

      <div className="flex items-center gap-3">
        <button
          className="md:hidden p-2"
          onClick={() => setShowMenu(!showMenu)}
        >
          <Menu className="w-5 h-5 text-ink" />
        </button>

        {user && (
          <div className="flex items-center gap-3">
            <Link href="/settings" className="flex items-center gap-2 hover:opacity-80 transition-opacity">
              <AvatarCircle name={user.name} src={user.avatar_url} size="sm" />
              <span className="hidden sm:block text-body-sm text-ink">{user.name}</span>
            </Link>
            <Link
              href="/settings"
              className="hidden sm:flex p-2 text-muted hover:text-ink transition-colors"
              title="Configurações"
            >
              <Settings className="w-4 h-4" />
            </Link>
            <button
              onClick={logout}
              className="p-2 text-muted hover:text-ink transition-colors"
              title="Sair"
            >
              <LogOut className="w-4 h-4" />
            </button>
          </div>
        )}
      </div>

      {showMenu && (
        <div className="md:hidden absolute top-16 left-0 right-0 bg-canvas border-b border-hairline-soft p-4 shadow-subtle-md">
          <nav className="flex flex-col gap-2">
            <MobileNavLink href="/dashboard" icon={<LayoutDashboard className="w-4 h-4" />}>Dashboard</MobileNavLink>
            <MobileNavLink href="/projects" icon={<FolderKanban className="w-4 h-4" />}>Projetos</MobileNavLink>
            <MobileNavLink href="/reports" icon={<BarChart3 className="w-4 h-4" />}>Relatórios</MobileNavLink>
            <MobileNavLink href="/settings" icon={<Settings className="w-4 h-4" />}>Configurações</MobileNavLink>
          </nav>
        </div>
      )}
    </header>
  );
}

function NavLink({ href, children, icon }: { href: string; children: React.ReactNode; icon: React.ReactNode }) {
  return (
    <Link
      href={href}
      className="flex items-center gap-2 px-3 py-2 text-nav-link font-medium text-muted hover:text-ink transition-colors"
    >
      {icon}
      {children}
    </Link>
  );
}

function MobileNavLink({ href, children, icon }: { href: string; children: React.ReactNode; icon: React.ReactNode }) {
  return (
    <Link
      href={href}
      className="flex items-center gap-3 px-3 py-2.5 text-body-sm font-medium text-muted hover:text-ink hover:bg-surface-soft rounded-md transition-colors"
    >
      {icon}
      {children}
    </Link>
  );
}
