'use client';

import { useEffect, useState } from 'react';
import { DashboardLayout } from '@/components/layout/dashboard-layout';
import { reportsApi } from '@/lib/api/reports';
import { DashboardData } from '@/types/report';
import { LoadingSpinner } from '@/components/shared/loading-spinner';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';

export default function ReportsPage() {
  const [data, setData] = useState<DashboardData | null>(null);
  const [loading, setLoading] = useState(true);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
    reportsApi.dashboard()
      .then((res) => setData(res.data.data))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <DashboardLayout>
        <LoadingSpinner className="py-20" />
      </DashboardLayout>
    );
  }

  if (!data) {
    return (
      <DashboardLayout>
        <div className="text-center py-20 text-muted">Erro ao carregar relatórios</div>
      </DashboardLayout>
    );
  }

  const COLORS = ['#fb923c', '#ec4899', '#8b5cf6', '#34d399'];

  const projectPriorityData = Object.entries(data.projects_by_priority || {}).map(([name, value]) => ({
    name: priorityLabel(name),
    value,
  }));

  return (
    <DashboardLayout>
      <h1 className="font-display text-display-md text-ink mb-8">Relatórios</h1>

      {mounted && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="bg-canvas rounded-lg border border-hairline p-6">
            <h3 className="font-body text-title-sm font-semibold text-ink mb-4">Projetos por Prioridade</h3>
            {projectPriorityData.length === 0 ? (
              <div className="h-[250px] flex items-center justify-center text-muted text-body-sm">
                Nenhum dado disponível
              </div>
            ) : (
              <ResponsiveContainer width="100%" height={250}>
                <BarChart data={projectPriorityData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="var(--hairline)" />
                  <XAxis dataKey="name" tick={{ fontSize: 12, fill: 'var(--muted)' }} />
                  <YAxis tick={{ fontSize: 12, fill: 'var(--muted)' }} />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: 'var(--surface-card)',
                      border: '1px solid var(--hairline)',
                      borderRadius: '8px',
                      color: 'var(--ink)',
                    }}
                  />
                  <Bar dataKey="value" fill="var(--brand-accent)" radius={[4, 4, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            )}
          </div>

          <div className="bg-canvas rounded-lg border border-hairline p-6">
            <h3 className="font-body text-title-sm font-semibold text-ink mb-4">Tarefas por Status</h3>
            {(data.tasks_by_status && Object.keys(data.tasks_by_status).length > 0) ? (
              <ResponsiveContainer width="100%" height={250}>
                <PieChart>
                  <Pie
                    data={Object.entries(data.tasks_by_status || {}).map(([name, value]) => ({
                      name: taskStatusLabel(name),
                      value,
                    }))}
                    cx="50%"
                    cy="50%"
                    outerRadius={80}
                    dataKey="value"
                  >
                    {Object.entries(data.tasks_by_status || {}).map((_, index) => (
                      <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip
                    contentStyle={{
                      backgroundColor: 'var(--surface-card)',
                      border: '1px solid var(--hairline)',
                      borderRadius: '8px',
                      color: 'var(--ink)',
                    }}
                  />
                </PieChart>
              </ResponsiveContainer>
            ) : (
              <div className="h-[250px] flex items-center justify-center text-muted text-body-sm">
                Nenhum dado disponível
              </div>
            )}
          </div>
        </div>
      )}

      <div className="mt-6 bg-surface-card rounded-lg p-6">
        <h3 className="font-body text-title-sm font-semibold text-ink mb-4">Resumo</h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <StatBox label="Total de Projetos" value={data.total_projects} />
          <StatBox label="Projetos Ativos" value={data.active_projects} />
          <StatBox label="Projetos Concluídos" value={data.completed_projects} />
          <StatBox label="Total de Tarefas" value={data.total_tasks} />
        </div>
      </div>
    </DashboardLayout>
  );
}

function StatBox({ label, value }: { label: string; value: number }) {
  return (
    <div className="bg-canvas rounded-md border border-hairline p-4 text-center">
      <p className="font-display text-display-sm text-ink">{value}</p>
      <p className="text-caption text-muted mt-1">{label}</p>
    </div>
  );
}

function priorityLabel(priority: string): string {
  const labels: Record<string, string> = {
    low: 'Baixa',
    medium: 'Média',
    high: 'Alta',
    critical: 'Crítica',
  };
  return labels[priority] || priority;
}

function taskStatusLabel(status: string): string {
  const labels: Record<string, string> = {
    todo: 'A Fazer',
    in_progress: 'Em Progresso',
    review: 'Revisão',
    done: 'Concluída',
  };
  return labels[status] || status;
}
