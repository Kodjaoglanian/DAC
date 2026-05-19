'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { DashboardLayout } from '@/components/layout/dashboard-layout';
import { projectsApi } from '@/lib/api/projects';
import { TextInput } from '@/components/shared/text-input';
import { ButtonPrimary } from '@/components/shared/button-primary';
import { ButtonSecondary } from '@/components/shared/button-secondary';
import { toast } from 'sonner';

export default function NewProjectPage() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [form, setForm] = useState({
    name: '',
    description: '',
    priority: 'medium' as 'low' | 'medium' | 'high' | 'critical',
    start_date: '',
    end_date: '',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    try {
      await projectsApi.create(form);
      toast.success('Projeto criado com sucesso!');
      router.push('/projects');
    } catch (error: any) {
      toast.error(error.response?.data?.error?.message || 'Erro ao criar projeto');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <DashboardLayout>
      <h1 className="font-display text-display-md text-ink mb-8">Novo Projeto</h1>

      <form onSubmit={handleSubmit} className="max-w-2xl bg-surface-card rounded-lg p-8 space-y-6">
        <div>
          <label className="block text-body-sm font-medium text-ink mb-1.5">Nome</label>
          <TextInput
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            placeholder="Nome do projeto"
            required
          />
        </div>

        <div>
          <label className="block text-body-sm font-medium text-ink mb-1.5">Descrição</label>
          <textarea
            value={form.description}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
            placeholder="Descreva o projeto..."
            rows={4}
            className="w-full bg-canvas text-ink font-body text-base rounded-md px-3.5 py-2.5 border border-hairline focus:border-ink focus:outline-none placeholder:text-muted transition-colors resize-none"
            required
          />
        </div>

        <div>
          <label className="block text-body-sm font-medium text-ink mb-1.5">Prioridade</label>
          <select
            value={form.priority}
            onChange={(e) => setForm({ ...form, priority: e.target.value as any })}
            className="w-full h-10 bg-canvas text-ink font-body text-base rounded-md px-3.5 border border-hairline focus:border-ink focus:outline-none"
          >
            <option value="low">Baixa</option>
            <option value="medium">Média</option>
            <option value="high">Alta</option>
            <option value="critical">Crítica</option>
          </select>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label className="block text-body-sm font-medium text-ink mb-1.5">Data de início</label>
            <TextInput
              type="date"
              value={form.start_date}
              onChange={(e) => setForm({ ...form, start_date: e.target.value })}
            />
          </div>
          <div>
            <label className="block text-body-sm font-medium text-ink mb-1.5">Data de término</label>
            <TextInput
              type="date"
              value={form.end_date}
              onChange={(e) => setForm({ ...form, end_date: e.target.value })}
            />
          </div>
        </div>

        <div className="flex gap-3 pt-4">
          <ButtonPrimary type="submit" isLoading={isLoading}>
            Criar Projeto
          </ButtonPrimary>
          <ButtonSecondary type="button" onClick={() => router.push('/projects')}>
            Cancelar
          </ButtonSecondary>
        </div>
      </form>
    </DashboardLayout>
  );
}
