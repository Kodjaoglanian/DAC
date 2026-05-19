'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/providers/auth-provider';
import { TextInput } from '@/components/shared/text-input';
import { ButtonPrimary } from '@/components/shared/button-primary';
import { LoadingSpinner } from '@/components/shared/loading-spinner';
import { toast } from 'sonner';

export default function LoginPage() {
  const { login, isLoading } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await login(email, password);
      toast.success('Login realizado com sucesso!');
    } catch (error: any) {
      toast.error(error.response?.data?.error?.message || 'Erro ao fazer login');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-canvas px-4">
      <div className="w-full max-w-md bg-surface-card rounded-lg p-8">
        <h1 className="font-display text-display-sm text-ink mb-1">Entrar na sua conta</h1>
        <p className="text-muted text-body-sm mb-8">ProjectTracker</p>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-body-sm font-medium text-ink mb-1.5">E-mail</label>
            <TextInput
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="seu@email.com"
              required
            />
          </div>

          <div>
            <label className="block text-body-sm font-medium text-ink mb-1.5">Senha</label>
            <TextInput
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              required
            />
          </div>

          <ButtonPrimary type="submit" className="w-full" isLoading={isLoading}>
            Entrar
          </ButtonPrimary>
        </form>

        <p className="text-center text-body-sm text-muted mt-6">
          Não tem uma conta?{' '}
          <Link href="/register" className="text-ink font-medium hover:underline">
            Criar conta
          </Link>
        </p>
      </div>
    </div>
  );
}
