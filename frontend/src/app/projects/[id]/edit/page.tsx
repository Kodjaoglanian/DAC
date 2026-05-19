'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { DashboardLayout } from '@/components/layout/dashboard-layout';
import { projectsApi } from '@/lib/api/projects';
import { TextInput } from '@/components/shared/text-input';
import { ButtonPrimary } from '@/components/shared/button-primary';
import { ButtonSecondary } from '@/components/shared/button-secondary';
import { toast } from 'sonner';
import { LoadingSpinner } from '@/components/shared/loading-spinner';

export default function EditProjectPage() {
  const params = useParams();
  const router = useRouter();
  const projectId = params.id as string;

  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [form, setForm] = useState({
    name: '',
    description: '',
    status: 'planning' as 'planning' | 'in_progress' | 'completed' | 'cancelled',
    priority: 'medium' as 'low' | 'medium' | 'high' | 'critical',
    start_date: '',
    end_date: '',
  });

  useEffect(() => {
    projectsApi.getById(projectId).then((res) => {
      const p = res.data.data;
      setForm({
        name: p.name,
        description: p.description,
        status: p.status,
        priority: p.priority,
        start_date: p.start_date ? p.start_date.slice(0, 10) : '',
        end_date: p.end_date ? p.end_date.slice(0, 10) : '',
      });
      setIsLoading(false);
    }).catch(() => {
      toast.error('Erro ao carregar projeto');
      router.push('/projects');
    });
  }, [projectId]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSaving(true);
    try {
      await projectsApi.update(projectId, form);
      toast.success('Projeto atualizado com sucesso!');
      router.push(`/projects/${projectId}`);
    } catch (error: any) {
      toast.error(error.response?.data?.error?.message || 'Erro ao atualizar projeto');
    } finally {
      setIsSaving(false);
    }
  };

  if (isLoading) {
    return (
      <DashboardLayout>
        <LoadingSpinner className="py-20" />
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <h1 className="font-display text-display-md text-ink mb-8">Editar Projeto</h1>

      <form onSubmit={handleSubmit} className="max-w-2xl bg-surface-card rounded-lg p-8 space-y-6">
        <div>
          <label className="block text-body-sm font-medium text-ink mb-1.5">Nome</label>
          <TextInput
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            required
          />
        </div>

        <div>
          <label className="block text-body-sm font-medium text-ink mb-1.5">Descrição</label>
          <textarea
            value={form.description}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
            rows={4}
            className="w-full bg-canvas text-ink font-body text-base rounded-md px-3.5 py-2.5 border border-hairline focus:border-ink focus:outline-none placeholder:text-muted transition-colors resize-none"
            required
          />
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label className="block text-body-sm font-medium text-ink mb-1.5">Status</label>
            <select
              value={form.status}
              onChange={(e) => setForm({ ...form, status: e.target.value as typeof form.status })}
              className="w-full h-10 bg-canvas text-ink font-body text-base rounded-md px-3.5 border border-hairline focus:border-ink focus:outline-none"
            >
              <option value="planning">Planejamento</option>
              <option value="in_progress">Em Andamento</option>
              <option value="completed">Concluído</option>
              <option value="cancelled">Cancelado</option>
            </select>
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
          <ButtonPrimary type="submit" isLoading={isSaving}>
            Salvar
          </ButtonPrimary>
          <ButtonSecondary type="button" onClick={() => router.push(`/projects/${projectId}`)}>
            Cancelar
          </ButtonSecondary>
        </div>
      </form>
    </DashboardLayout>
  );
}
