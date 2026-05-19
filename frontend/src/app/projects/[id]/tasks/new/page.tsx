'use client';

import { useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { DashboardLayout } from '@/components/layout/dashboard-layout';
import { tasksApi } from '@/lib/api/tasks';
import { TextInput } from '@/components/shared/text-input';
import { ButtonPrimary } from '@/components/shared/button-primary';
import { ButtonSecondary } from '@/components/shared/button-secondary';
import { toast } from 'sonner';

export default function NewTaskPage() {
  const params = useParams();
  const router = useRouter();
  const projectId = params.id as string;

  const [isLoading, setIsLoading] = useState(false);
  const [form, setForm] = useState({
    title: '',
    description: '',
    priority: 'medium' as 'low' | 'medium' | 'high' | 'critical',
    due_date: '',
    assigned_to: '',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    try {
      await tasksApi.create(projectId, form);
      toast.success('Tarefa criada com sucesso!');
      router.push(`/projects/${projectId}`);
    } catch (error: any) {
      toast.error(error.response?.data?.error?.message || 'Erro ao criar tarefa');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <DashboardLayout>
      <h1 className="font-display text-display-md text-ink mb-8">Nova Tarefa</h1>

      <form onSubmit={handleSubmit} className="max-w-2xl bg-surface-card rounded-lg p-8 space-y-6">
        <div>
          <label className="block text-body-sm font-medium text-ink mb-1.5">Título</label>
          <TextInput
            value={form.title}
            onChange={(e) => setForm({ ...form, title: e.target.value })}
            placeholder="Título da tarefa"
            required
          />
        </div>

        <div>
          <label className="block text-body-sm font-medium text-ink mb-1.5">Descrição</label>
          <textarea
            value={form.description}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
            placeholder="Descreva a tarefa..."
            rows={3}
            className="w-full bg-canvas text-ink font-body text-base rounded-md px-3.5 py-2.5 border border-hairline focus:border-ink focus:outline-none placeholder:text-muted transition-colors resize-none"
          />
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
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
          <div>
            <label className="block text-body-sm font-medium text-ink mb-1.5">Data de vencimento</label>
            <TextInput
              type="date"
              value={form.due_date}
              onChange={(e) => setForm({ ...form, due_date: e.target.value })}
            />
          </div>
        </div>

        <div className="flex gap-3 pt-4">
          <ButtonPrimary type="submit" isLoading={isLoading}>
            Criar Tarefa
          </ButtonPrimary>
          <ButtonSecondary type="button" onClick={() => router.push(`/projects/${projectId}`)}>
            Cancelar
          </ButtonSecondary>
        </div>
      </form>
    </DashboardLayout>
  );
}
