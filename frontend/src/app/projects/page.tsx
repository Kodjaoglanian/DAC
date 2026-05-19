'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { DashboardLayout } from '@/components/layout/dashboard-layout';
import { projectsApi } from '@/lib/api/projects';
import { Project, ProjectStatus, Priority } from '@/types/project';
import { LoadingSpinner } from '@/components/shared/loading-spinner';
import { BadgePill } from '@/components/shared/badge-pill';
import { ButtonPrimary } from '@/components/shared/button-primary';
import { TextInput } from '@/components/shared/text-input';
import { Plus, Search } from 'lucide-react';

export default function ProjectsPage() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState<ProjectStatus | ''>('');
  const [priorityFilter, setPriorityFilter] = useState<Priority | ''>('');

  useEffect(() => {
    loadProjects();
  }, [statusFilter, priorityFilter]);

  const loadProjects = async () => {
    setLoading(true);
    try {
      const res = await projectsApi.list({
        status: statusFilter || undefined,
        priority: priorityFilter || undefined,
        search: search || undefined,
      });
      setProjects(res.data.data?.projects || []);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    loadProjects();
  };

  const filteredProjects = projects.filter((p) =>
    search ? p.name.toLowerCase().includes(search.toLowerCase()) : true
  );

  return (
    <DashboardLayout>
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8">
        <h1 className="font-display text-display-md text-ink">Projetos</h1>
        <Link href="/projects/new">
          <ButtonPrimary>
            <Plus className="w-4 h-4 mr-1.5" />
            Novo Projeto
          </ButtonPrimary>
        </Link>
      </div>

      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-3 mb-6">
        <form onSubmit={handleSearch} className="flex gap-2 flex-1">
          <TextInput
            placeholder="Buscar projetos..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="flex-1"
          />
          <button type="submit" className="p-2.5 bg-surface-card rounded-md text-ink">
            <Search className="w-4 h-4" />
          </button>
        </form>

        <select
          value={statusFilter}
          onChange={(e) => setStatusFilter(e.target.value as ProjectStatus)}
          className="h-10 px-3 bg-canvas border border-hairline rounded-md text-body-sm text-ink"
        >
          <option value="">Todos os status</option>
          <option value="planning">Planejamento</option>
          <option value="in_progress">Em Andamento</option>
          <option value="completed">Concluído</option>
          <option value="cancelled">Cancelado</option>
        </select>

        <select
          value={priorityFilter}
          onChange={(e) => setPriorityFilter(e.target.value as Priority)}
          className="h-10 px-3 bg-canvas border border-hairline rounded-md text-body-sm text-ink"
        >
          <option value="">Todas as prioridades</option>
          <option value="low">Baixa</option>
          <option value="medium">Média</option>
          <option value="high">Alta</option>
          <option value="critical">Crítica</option>
        </select>
      </div>

      {loading ? (
        <LoadingSpinner className="py-20" />
      ) : filteredProjects.length === 0 ? (
        <div className="text-center py-20 bg-surface-card rounded-lg">
          <p className="text-muted text-body-md">Nenhum projeto encontrado</p>
          <Link href="/projects/new" className="text-ink font-medium hover:underline mt-2 inline-block">
            Criar primeiro projeto
          </Link>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredProjects.map((project) => (
            <ProjectCard key={project.id} project={project} />
          ))}
        </div>
      )}
    </DashboardLayout>
  );
}

function ProjectCard({ project }: { project: Project }) {
  const statusVariant = {
    planning: 'default' as const,
    in_progress: 'orange' as const,
    completed: 'emerald' as const,
    cancelled: 'pink' as const,
  };

  const priorityVariant = {
    low: 'default' as const,
    medium: 'default' as const,
    high: 'orange' as const,
    critical: 'pink' as const,
  };

  return (
    <Link href={`/projects/${project.id}`}>
      <div className="bg-surface-card rounded-lg p-8 hover:shadow-subtle-md transition-shadow cursor-pointer h-full flex flex-col">
        <div className="flex items-center gap-2 mb-3">
          <BadgePill label={statusLabel(project.status)} variant={statusVariant[project.status]} />
        </div>

        <h3 className="font-body text-title-md font-semibold text-ink mb-2">{project.name}</h3>
        <p className="text-body text-body-sm mb-4 line-clamp-2 flex-1">{project.description}</p>

        <div className="flex items-center gap-2 mb-4">
          <BadgePill label={priorityLabel(project.priority)} variant={priorityVariant[project.priority]} />
        </div>

        <div className="flex items-center justify-between text-caption text-muted">
          <span>{project.tasks_count} tarefas</span>
          <span>{project.created_by.name}</span>
        </div>
      </div>
    </Link>
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

function priorityLabel(priority: string): string {
  const labels: Record<string, string> = {
    low: 'Baixa',
    medium: 'Média',
    high: 'Alta',
    critical: 'Crítica',
  };
  return labels[priority] || priority;
}
