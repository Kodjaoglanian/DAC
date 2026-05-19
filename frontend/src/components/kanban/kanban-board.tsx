'use client';

import { useState } from 'react';
import { Task } from '@/types/task';
import { tasksApi } from '@/lib/api/tasks';
import { BadgePill } from '@/components/shared/badge-pill';
import { toast } from 'sonner';
import { Calendar, User } from 'lucide-react';

interface KanbanBoardProps {
  projectId: string;
  tasks: Task[];
  onUpdate: () => void;
}

const COLUMNS = [
  { id: 'todo', label: 'A Fazer', variant: 'default' as const },
  { id: 'in_progress', label: 'Em Progresso', variant: 'orange' as const },
  { id: 'review', label: 'Revisão', variant: 'violet' as const },
  { id: 'done', label: 'Concluída', variant: 'emerald' as const },
];

export function KanbanBoard({ projectId, tasks, onUpdate }: KanbanBoardProps) {
  const [draggedTask, setDraggedTask] = useState<string | null>(null);
  const [dragOverColumn, setDragOverColumn] = useState<string | null>(null);

  const handleDragStart = (taskId: string) => {
    setDraggedTask(taskId);
  };

  const handleDragOver = (e: React.DragEvent, columnId: string) => {
    e.preventDefault();
    setDragOverColumn(columnId);
  };

  const handleDragLeave = () => {
    setDragOverColumn(null);
  };

  const handleDrop = async (e: React.DragEvent, newStatus: string) => {
    e.preventDefault();
    setDragOverColumn(null);

    if (!draggedTask) return;

    const task = tasks.find((t) => t.id === draggedTask);
    if (!task || task.status === newStatus) {
      setDraggedTask(null);
      return;
    }

    try {
      await tasksApi.updateStatus(draggedTask, newStatus);
      toast.success('Status atualizado');
      onUpdate();
    } catch {
      toast.error('Erro ao mover tarefa');
    } finally {
      setDraggedTask(null);
    }
  };

  const getTasksByColumn = (status: string) => tasks.filter((t) => t.status === status);

  return (
    <div className="flex gap-4 overflow-x-auto pb-4 min-h-[400px]">
      {COLUMNS.map((column) => {
        const columnTasks = getTasksByColumn(column.id);
        const isDragOver = dragOverColumn === column.id;

        return (
          <div
            key={column.id}
            className={`flex-shrink-0 w-72 bg-surface-soft rounded-lg border border-hairline transition-colors ${
              isDragOver ? 'bg-surface-card border-brand-accent' : ''
            }`}
            onDragOver={(e) => handleDragOver(e, column.id)}
            onDragLeave={handleDragLeave}
            onDrop={(e) => handleDrop(e, column.id)}
          >
            <div className="p-3 border-b border-hairline flex items-center justify-between">
              <div className="flex items-center gap-2">
                <BadgePill label={column.label} variant={column.variant} />
              </div>
              <span className="text-caption text-muted font-medium">{columnTasks.length}</span>
            </div>

            <div className="p-3 space-y-3 min-h-[200px]">
              {columnTasks.map((task) => (
                <div
                  key={task.id}
                  draggable
                  onDragStart={() => handleDragStart(task.id)}
                  className={`bg-canvas rounded-md border border-hairline p-3 cursor-move hover:shadow-subtle transition-shadow ${
                    draggedTask === task.id ? 'opacity-50' : ''
                  }`}
                >
                  <h4 className="text-body-sm font-medium text-ink mb-2">{task.title}</h4>
                  {task.description && (
                    <p className="text-caption text-muted mb-2 line-clamp-2">{task.description}</p>
                  )}
                  <div className="flex items-center gap-2 flex-wrap">
                    <BadgePill
                      label={priorityLabel(task.priority)}
                      variant={task.priority === 'high' || task.priority === 'critical' ? 'pink' : 'default'}
                    />
                    {task.due_date && (
                      <span className="flex items-center gap-1 text-caption text-muted">
                        <Calendar className="w-3 h-3" />
                        {formatDate(task.due_date)}
                      </span>
                    )}
                    {task.assigned_to && (
                      <span className="flex items-center gap-1 text-caption text-muted">
                        <User className="w-3 h-3" />
                        {task.assigned_to.name}
                      </span>
                    )}
                  </div>
                </div>
              ))}
              {columnTasks.length === 0 && !isDragOver && (
                <div className="text-center py-8 text-muted text-caption">Arraste tarefas aqui</div>
              )}
            </div>
          </div>
        );
      })}
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

function formatDate(date: string): string {
  return new Date(date).toLocaleDateString('pt-BR', { day: '2-digit', month: '2-digit' });
}
