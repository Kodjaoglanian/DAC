'use client';

import { Task } from '@/types/task';
import { BadgePill } from '@/components/shared/badge-pill';
import { CalendarDays, ChevronLeft, ChevronRight } from 'lucide-react';
import { useState } from 'react';

interface ProjectCalendarProps {
  tasks: Task[];
}

export function ProjectCalendar({ tasks }: ProjectCalendarProps) {
  const [currentDate, setCurrentDate] = useState(new Date());

  const year = currentDate.getFullYear();
  const month = currentDate.getMonth();

  const firstDay = new Date(year, month, 1);
  const lastDay = new Date(year, month + 1, 0);
  const daysInMonth = lastDay.getDate();
  const startDayOfWeek = firstDay.getDay();

  const monthName = currentDate.toLocaleDateString('pt-BR', { month: 'long', year: 'numeric' });

  const prevMonth = () => setCurrentDate(new Date(year, month - 1, 1));
  const nextMonth = () => setCurrentDate(new Date(year, month + 1, 1));

  const tasksWithDueDate = tasks.filter((t) => t.due_date);

  const getTasksForDay = (day: number) => {
    const dateStr = new Date(year, month, day).toISOString().split('T')[0];
    return tasksWithDueDate.filter((t) => t.due_date && t.due_date.startsWith(dateStr));
  };

  const days = [];
  for (let i = 0; i < startDayOfWeek; i++) {
    days.push(<div key={`empty-${i}`} className="h-24 bg-surface-soft/50 rounded-md" />);
  }

  for (let day = 1; day <= daysInMonth; day++) {
    const dayTasks = getTasksForDay(day);
    const isToday = new Date().toDateString() === new Date(year, month, day).toDateString();

    days.push(
      <div
        key={day}
        className={`h-24 rounded-md border p-1.5 overflow-hidden ${
          isToday ? 'bg-brand-accent/5 border-brand-accent' : 'bg-canvas border-hairline'
        }`}
      >
        <span className={`text-caption font-medium ${isToday ? 'text-brand-accent' : 'text-muted'}`}>
          {day}
        </span>
        <div className="mt-1 space-y-1">
          {dayTasks.slice(0, 2).map((task) => (
            <div
              key={task.id}
              className="text-[10px] leading-tight truncate px-1.5 py-0.5 rounded bg-surface-card text-ink"
              title={task.title}
            >
              {task.title}
            </div>
          ))}
          {dayTasks.length > 2 && (
            <div className="text-[10px] text-muted px-1.5">+{dayTasks.length - 2} mais</div>
          )}
        </div>
      </div>
    );
  }

  const weekDays = ['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb'];

  return (
    <div className="bg-canvas rounded-lg border border-hairline p-6">
      <div className="flex items-center justify-between mb-6">
        <h3 className="font-body text-title-sm font-semibold text-ink flex items-center gap-2">
          <CalendarDays className="w-5 h-5" />
          Calendário
        </h3>
        <div className="flex items-center gap-2">
          <button onClick={prevMonth} className="p-1.5 hover:bg-surface-soft rounded-md transition-colors">
            <ChevronLeft className="w-4 h-4 text-muted" />
          </button>
          <span className="text-body-sm font-medium text-ink capitalize min-w-[140px] text-center">{monthName}</span>
          <button onClick={nextMonth} className="p-1.5 hover:bg-surface-soft rounded-md transition-colors">
            <ChevronRight className="w-4 h-4 text-muted" />
          </button>
        </div>
      </div>

      <div className="grid grid-cols-7 gap-2">
        {weekDays.map((wd) => (
          <div key={wd} className="text-center text-caption text-muted font-medium py-2">
            {wd}
          </div>
        ))}
        {days}
      </div>

      <div className="mt-4 flex items-center gap-4 text-caption text-muted">
        <div className="flex items-center gap-1.5">
          <div className="w-2 h-2 rounded-full bg-brand-accent" />
          Com prazo
        </div>
        <div className="flex items-center gap-1.5">
          <div className="w-2 h-2 rounded-full bg-badge-orange" />
          Alta prioridade
        </div>
      </div>
    </div>
  );
}
