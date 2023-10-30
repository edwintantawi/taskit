import React, { useState } from 'react';
import { format } from 'date-fns';
import {
  CalendarIcon,
  TrashIcon,
  PencilSquareIcon,
} from '@heroicons/react/24/outline';

import { TaskCheck } from '~/task/components';
import { Task } from '~/task/entity';
import { useDeleteTask, useUpdateTask } from '~/task/hooks';
import { TaskEditor } from '~/task/containers';

export function TaskItem({
  id,
  content,
  description,
  due_date,
  is_completed,
}: Task) {
  const [isEditorOpen, setIsEditorOpen] = useState(false);
  const { mutate: deleteTask } = useDeleteTask(id);
  const { mutate: updateTask, isLoading } = useUpdateTask(id);

  const handleDeleteTask = () => deleteTask();
  const handleOpenEditor = () => setIsEditorOpen(true);
  const handleCancelEdit = () => setIsEditorOpen(false);

  const handleChangeCheck = () =>
    updateTask({
      id,
      content,
      description,
      due_date,
      is_completed: !is_completed,
    });

  if (isEditorOpen) {
    return (
      <TaskEditor
        onCancel={handleCancelEdit}
        onSubmit={(data) => updateTask(data, { onSuccess: handleCancelEdit })}
        isLoading={isLoading}
        submitTitle="Update task"
        submitLoadingTitle="Updating new task, please wait..."
        data={{ id, content, description, due_date, is_completed }}
      />
    );
  }

  return (
    <li className="group flex gap-3 border-b py-4">
      <TaskCheck onChange={handleChangeCheck} checked={is_completed} />
      <div className="relative w-[calc(100%-32px)] flex-1">
        <h3 className="mb-1 w-full break-words text-sm leading-[1.4rem]">
          {content}
        </h3>
        <p className="mb-2 text-xs text-gray-600">
          {description || 'no description'}
        </p>
        {due_date && (
          <p className="flex items-center gap-1 text-xs text-gray-500">
            <CalendarIcon className="h-4 w-4" />{' '}
            {format(new Date(due_date), 'dd MMMM yyyy')}
          </p>
        )}

        <div className="invisible absolute top-0 right-0 flex h-10 text-sm text-gray-500 group-hover:visible">
          <div className="w-10  bg-gradient-to-l from-white to-transparent" />
          <div className="flex items-center gap-1 bg-white pl-3">
            <button
              className="rounded-md p-1 hover:bg-gray-100"
              onClick={handleOpenEditor}
            >
              <PencilSquareIcon className="h-5 w-5" />
            </button>
            <button
              className="rounded-md p-1 text-red-500 hover:bg-gray-100"
              onClick={handleDeleteTask}
            >
              <TrashIcon className="h-5 w-5" />
            </button>
          </div>
        </div>
      </div>
    </li>
  );
}
