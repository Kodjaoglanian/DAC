'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { DashboardLayout } from '@/components/layout/dashboard-layout';
import { projectsApi } from '@/lib/api/projects';
import { tasksApi } from '@/lib/api/tasks';
import { membersApi } from '@/lib/api/members';
import { Project } from '@/types/project';
import { Task } from '@/types/task';
import { Member } from '@/types/member';
import { LoadingSpinner } from '@/components/shared/loading-spinner';
import { BadgePill } from '@/components/shared/badge-pill';
import { ButtonPrimary } from '@/components/shared/button-primary';
import { ButtonSecondary } from '@/components/shared/button-secondary';
import { AvatarCircle } from '@/components/shared/avatar-circle';
import { toast } from 'sonner';
import { Edit, Trash2, Plus, List, LayoutGrid } from 'lucide-react';
import Link from 'next/link';
import { KanbanBoard } from '@/components/kanban/kanban-board';
import { ProjectCalendar } from '@/components/project/project-calendar';

export default function ProjectDetailPage() {
  const params = useParams();
  const router = useRouter();
  const projectId = params.id as string;

  const [project, setProject] = useState<Project | null>(null);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [members, setMembers] = useState<Member[]>([]);
  const [activeTab, setActiveTab] = useState<'tasks' | 'members' | 'calendar'>('tasks');
  const [taskView, setTaskView] = useState<'list' | 'kanban'>('kanban');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, [projectId]);

  const loadData = async () => {
    setLoading(true);
    try {
      const [projectRes, tasksRes, membersRes] = await Promise.all([
        projectsApi.getById(projectId),
        tasksApi.list(projectId),
        membersApi.list(projectId),
      ]);
      setProject(projectRes.data.data);
      setTasks((tasksRes.data.data as any)?.tasks || []);
      setMembers(membersRes.data.data || []);
    } catch {
      toast.error('Erro ao carregar projeto');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!confirm('Tem certeza que deseja excluir este projeto?')) return;
    try {
      await projectsApi.delete(projectId);
      toast.success('Projeto excluído');
      router.push('/projects');
    } catch {
      toast.error('Erro ao excluir projeto');
    }
  };

  const handleUpdateTaskStatus = async (taskId: string, status: string) => {
    try {
      await tasksApi.updateStatus(taskId, status);
      toast.success('Status atualizado');
      loadData();
    } catch {
      toast.error('Erro ao atualizar status');
    }
  };

  if (loading) {
    return (
      <DashboardLayout>
        <LoadingSpinner className="py-20" />
      </DashboardLayout>
    );
  }

  if (!project) {
    return (
      <DashboardLayout>
        <div className="text-center py-20 text-muted">Projeto não encontrado</div>
      </DashboardLayout>
    );
  }

  const statusVariant = {
    planning: 'default' as const,
    in_progress: 'orange' as const,
    completed: 'emerald' as const,
    cancelled: 'pink' as const,
  };

  const taskStatusVariant = {
    todo: 'default' as const,
    in_progress: 'orange' as const,
    review: 'violet' as const,
    done: 'emerald' as const,
  };

  return (
    <DashboardLayout>
      <div className="mb-8">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-4">
          <div>
            <h1 className="font-display text-display-md text-ink">{project.name}</h1>
            <p className="text-muted text-body-sm mt-1">{project.description}</p>
          </div>
          <div className="flex gap-2">
            <Link href={`/projects/${project.id}/edit`}>
              <ButtonSecondary>
                <Edit className="w-4 h-4 mr-1.5" />
                Editar
              </ButtonSecondary>
            </Link>
            <ButtonSecondary onClick={handleDelete} className="text-error hover:bg-error/10">
              <Trash2 className="w-4 h-4 mr-1.5" />
              Excluir
            </ButtonSecondary>
          </div>
        </div>

        <div className="flex items-center gap-2 mb-4">
          <BadgePill label={statusLabel(project.status)} variant={statusVariant[project.status]} />
          <BadgePill label={priorityLabel(project.priority)} variant="default" />
        </div>
      </div>

      {/* Tabs */}
      <div className="bg-surface-soft rounded-pill p-1.5 flex gap-1 mb-6 w-fit">
        <button
          onClick={() => setActiveTab('tasks')}
          className={`px-4 py-2 rounded-md font-body text-nav-link font-medium transition-all ${
            activeTab === 'tasks' ? 'bg-canvas text-ink shadow-subtle' : 'text-muted'
          }`}
        >
          Tarefas ({tasks.length})
        </button>
        <button
          onClick={() => setActiveTab('members')}
          className={`px-4 py-2 rounded-md font-body text-nav-link font-medium transition-all ${
            activeTab === 'members' ? 'bg-canvas text-ink shadow-subtle' : 'text-muted'
          }`}
        >
          Membros ({members.length})
        </button>
        <button
          onClick={() => setActiveTab('calendar')}
          className={`px-4 py-2 rounded-md font-body text-nav-link font-medium transition-all ${
            activeTab === 'calendar' ? 'bg-canvas text-ink shadow-subtle' : 'text-muted'
          }`}
        >
          Calendário
        </button>
      </div>

      {activeTab === 'tasks' ? (
        <div>
          <div className="flex justify-between items-center mb-4">
            <h2 className="font-body text-title-sm font-semibold text-ink">Tarefas</h2>
            <div className="flex items-center gap-2">
              <div className="bg-surface-soft rounded-pill p-1 flex gap-1">
                <button
                  onClick={() => setTaskView('kanban')}
                  className={`flex items-center gap-1.5 px-3 py-1.5 rounded-md text-body-sm font-medium transition-all ${
                    taskView === 'kanban' ? 'bg-canvas text-ink shadow-subtle' : 'text-muted'
                  }`}
                >
                  <LayoutGrid className="w-4 h-4" />
                  Kanban
                </button>
                <button
                  onClick={() => setTaskView('list')}
                  className={`flex items-center gap-1.5 px-3 py-1.5 rounded-md text-body-sm font-medium transition-all ${
                    taskView === 'list' ? 'bg-canvas text-ink shadow-subtle' : 'text-muted'
                  }`}
                >
                  <List className="w-4 h-4" />
                  Lista
                </button>
              </div>
              <Link href={`/projects/${projectId}/tasks/new`}>
                <ButtonPrimary>
                  <Plus className="w-4 h-4 mr-1.5" />
                  Nova Tarefa
                </ButtonPrimary>
              </Link>
            </div>
          </div>

          {tasks.length === 0 ? (
            <div className="text-center py-12 bg-surface-card rounded-lg">
              <p className="text-muted">Nenhuma tarefa ainda</p>
            </div>
          ) : taskView === 'kanban' ? (
            <KanbanBoard projectId={projectId} tasks={tasks} onUpdate={loadData} />
          ) : (
            <div className="space-y-3">
              {tasks.map((task) => (
                <div
                  key={task.id}
                  className="flex items-center justify-between p-4 bg-surface-card rounded-lg"
                >
                  <div>
                    <h3 className="font-body text-body-md font-medium text-ink">{task.title}</h3>
                    <div className="flex items-center gap-2 mt-1">
                      <BadgePill label={taskStatusLabel(task.status)} variant={taskStatusVariant[task.status]} />
                      {task.assigned_to && (
                        <span className="text-caption text-muted">{task.assigned_to.name}</span>
                      )}
                    </div>
                  </div>
                  <select
                    value={task.status}
                    onChange={(e) => handleUpdateTaskStatus(task.id, e.target.value)}
                    className="h-8 px-2 bg-canvas border border-hairline rounded-md text-caption text-ink"
                  >
                    <option value="todo">A Fazer</option>
                    <option value="in_progress">Em Progresso</option>
                    <option value="review">Revisão</option>
                    <option value="done">Concluída</option>
                  </select>
                </div>
              ))}
            </div>
          )}
        </div>
      ) : activeTab === 'members' ? (
        <div>
          <h2 className="font-body text-title-sm font-semibold text-ink mb-4">Membros</h2>
          {members.length === 0 ? (
            <div className="text-center py-12 bg-surface-card rounded-lg">
              <p className="text-muted">Nenhum membro</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
              {members.map((member) => (
                <div key={member.id} className="flex items-center gap-3 p-4 bg-surface-card rounded-lg">
                  <AvatarCircle name={member.user.name} src={member.user.avatar_url} />
                  <div>
                    <p className="text-ink font-medium text-body-sm">{member.user.name}</p>
                    <p className="text-caption text-muted capitalize">{member.role}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      ) : activeTab === 'calendar' ? (
        <ProjectCalendar tasks={tasks} />
      ) : null}
    </DashboardLayout>
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

function taskStatusLabel(status: string): string {
  const labels: Record<string, string> = {
    todo: 'A Fazer',
    in_progress: 'Em Progresso',
    review: 'Revisão',
    done: 'Concluída',
  };
  return labels[status] || status;
}
