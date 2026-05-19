'use client';

import { useEffect, useState } from 'react';
import { DashboardLayout } from '@/components/layout/dashboard-layout';
import { reportsApi } from '@/lib/api/reports';
import { DashboardData } from '@/types/report';
import { LoadingSpinner } from '@/components/shared/loading-spinner';
import { FolderKanban, Activity, CheckCircle2, ListTodo } from 'lucide-react';
import { PieChart, Pie, Cell, ResponsiveContainer, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip } from 'recharts';
import Link from 'next/link';

export default function DashboardPage() {
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
        <div className="text-center py-20 text-muted">Erro ao carregar dados</div>
      </DashboardLayout>
    );
  }

  const projectStatusData = Object.entries(data.projects_by_status || {}).map(([name, value]) => ({
    name: statusLabel(name),
    value,
  }));

  const taskStatusData = Object.entries(data.tasks_by_status || {}).map(([name, value]) => ({
    name: taskStatusLabel(name),
    value,
  }));

  const COLORS = ['#fb923c', '#ec4899', '#8b5cf6', '#34d399'];

  return (
    <DashboardLayout>
      <h1 className="font-display text-display-md text-ink mb-8">Dashboard</h1>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <StatCard icon={<FolderKanban className="w-5 h-5" />} label="Projetos" value={data.total_projects} />
        <StatCard icon={<Activity className="w-5 h-5" />} label="Ativos" value={data.active_projects} />
        <StatCard icon={<CheckCircle2 className="w-5 h-5" />} label="Concluídos" value={data.completed_projects} />
        <StatCard icon={<ListTodo className="w-5 h-5" />} label="Tarefas" value={data.total_tasks} />
      </div>

      {/* Charts */}
      {mounted && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          <div className="bg-canvas rounded-lg border border-hairline p-6">
            <h3 className="font-body text-title-sm font-semibold text-ink mb-4">Projetos por Status</h3>
            {projectStatusData.length === 0 ? (
              <div className="h-[250px] flex items-center justify-center text-muted text-body-sm">
                Nenhum dado disponível
              </div>
            ) : (
              <ResponsiveContainer width="100%" height={250}>
                <PieChart>
                  <Pie data={projectStatusData} cx="50%" cy="50%" outerRadius={80} dataKey="value">
                    {projectStatusData.map((_, index) => (
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
            )}
          </div>

          <div className="bg-canvas rounded-lg border border-hairline p-6">
            <h3 className="font-body text-title-sm font-semibold text-ink mb-4">Tarefas por Status</h3>
            {taskStatusData.length === 0 ? (
              <div className="h-[250px] flex items-center justify-center text-muted text-body-sm">
                Nenhum dado disponível
              </div>
            ) : (
              <ResponsiveContainer width="100%" height={250}>
                <BarChart data={taskStatusData}>
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
        </div>
      )}

      {/* Recent Projects */}
      <div className="bg-surface-card rounded-lg p-6">
        <h3 className="font-body text-title-sm font-semibold text-ink mb-4">Projetos Recentes</h3>
        {(data.recent_projects || []).length === 0 ? (
          <p className="text-muted text-body-sm">Nenhum projeto recente</p>
        ) : (
          <div className="space-y-3">
            {(data.recent_projects || []).map((project) => (
              <Link
                key={project.id}
                href={`/projects/${project.id}`}
                className="flex items-center justify-between p-3 bg-canvas rounded-md hover:bg-surface-soft transition-colors"
              >
                <div>
                  <p className="text-ink font-medium text-body-sm">{project.name}</p>
                  <p className="text-muted text-caption">{project.status}</p>
                </div>
                <span className="text-caption text-muted">{project.tasks_count} tarefas</span>
              </Link>
            ))}
          </div>
        )}
      </div>
    </DashboardLayout>
  );
}

function StatCard({ icon, label, value }: { icon: React.ReactNode; label: string; value: number }) {
  return (
    <div className="bg-canvas rounded-lg border border-hairline p-6">
      <div className="text-ink mb-3">{icon}</div>
      <p className="font-display text-display-sm text-ink">{value}</p>
      <p className="text-muted text-body-sm">{label}</p>
    </div>
  );
}

function statusLabel(status: string): string {
  const labels: Record<string, string> = {
    planning: 'Planejamento',
    in_progress: 'Em Andamento',
    completed: 'Concluído',
    cancelled: 'Cancelado',
  };
  return labels[status] || status;
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
