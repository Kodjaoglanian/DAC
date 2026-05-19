'use client';

import { useState } from 'react';
import { useAuth } from '@/providers/auth-provider';
import { DashboardLayout } from '@/components/layout/dashboard-layout';
import { TextInput } from '@/components/shared/text-input';
import { ButtonPrimary } from '@/components/shared/button-primary';
import { ButtonSecondary } from '@/components/shared/button-secondary';
import { AvatarCircle } from '@/components/shared/avatar-circle';
import { useTheme } from 'next-themes';
import { Sun, Moon, Monitor } from 'lucide-react';
import { toast } from 'sonner';

export default function SettingsPage() {
  const { user, logout } = useAuth();
  const { theme, setTheme } = useTheme();
  const [name, setName] = useState(user?.name || '');

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();
    toast.success('Perfil atualizado (mock)');
  };

  return (
    <DashboardLayout>
      <h1 className="font-display text-display-md text-ink mb-8">Configurações</h1>

      <div className="max-w-2xl space-y-6">
        {/* Profile */}
        <div className="bg-surface-card rounded-lg p-8">
          <h2 className="font-body text-title-md font-semibold text-ink mb-6">Perfil</h2>

          <div className="flex items-center gap-4 mb-6">
            <AvatarCircle name={user?.name || ''} src={user?.avatar_url} size="lg" />
            <div>
              <p className="text-ink font-medium">{user?.name}</p>
              <p className="text-muted text-body-sm">{user?.email}</p>
              <p className="text-caption text-muted capitalize">{user?.role}</p>
            </div>
          </div>

          <form onSubmit={handleSave} className="space-y-4">
            <div>
              <label className="block text-body-sm font-medium text-ink mb-1.5">Nome</label>
              <TextInput value={name} onChange={(e) => setName(e.target.value)} />
            </div>
            <div>
              <label className="block text-body-sm font-medium text-ink mb-1.5">E-mail</label>
              <TextInput value={user?.email || ''} disabled className="bg-surface-soft" />
            </div>
            <ButtonPrimary type="submit">Salvar Perfil</ButtonPrimary>
          </form>
        </div>

        {/* Theme */}
        <div className="bg-surface-card rounded-lg p-8">
          <h2 className="font-body text-title-md font-semibold text-ink mb-6">Aparência</h2>

          <div className="bg-surface-soft rounded-pill p-1.5 flex gap-1 w-fit">
            <button
              onClick={() => setTheme('light')}
              className={`flex items-center gap-2 px-4 py-2 rounded-md font-body text-nav-link font-medium transition-all ${
                theme === 'light' ? 'bg-canvas text-ink shadow-subtle' : 'text-muted'
              }`}
            >
              <Sun className="w-4 h-4" />
              Claro
            </button>
            <button
              onClick={() => setTheme('dark')}
              className={`flex items-center gap-2 px-4 py-2 rounded-md font-body text-nav-link font-medium transition-all ${
                theme === 'dark' ? 'bg-canvas text-ink shadow-subtle' : 'text-muted'
              }`}
            >
              <Moon className="w-4 h-4" />
              Escuro
            </button>
            <button
              onClick={() => setTheme('system')}
              className={`flex items-center gap-2 px-4 py-2 rounded-md font-body text-nav-link font-medium transition-all ${
                theme === 'system' ? 'bg-canvas text-ink shadow-subtle' : 'text-muted'
              }`}
            >
              <Monitor className="w-4 h-4" />
              Sistema
            </button>
          </div>
        </div>

        {/* Logout */}
        <div className="bg-surface-card rounded-lg p-8">
          <h2 className="font-body text-title-md font-semibold text-ink mb-4">Sessão</h2>
          <ButtonSecondary onClick={logout} className="text-error border-error/20 hover:bg-error/10">
            Sair da conta
          </ButtonSecondary>
        </div>
      </div>
    </DashboardLayout>
  );
}
