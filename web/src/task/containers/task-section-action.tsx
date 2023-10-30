import React, { useState } from 'react';
import { PlusIcon } from '@heroicons/react/24/outline';

import { Button } from '~/common/components';
import { useAddTask } from '~/task/hooks';
import { TaskEditor } from '~/task/containers';

export function TaskSectionAction() {
  const [isEditorOpen, setIsEditorOpen] = useState(false);
  const { mutate: addTask, isLoading } = useAddTask();

  const handleEditorOpen = () => setIsEditorOpen(true);
  const handleEditorClose = () => setIsEditorOpen(false);

  return (
    <li className="mt-4 space-y-4">
      <Button
        variants="normal"
        size="small"
        className="flex items-center gap-1 font-semibold"
        onClick={handleEditorOpen}
      >
        <PlusIcon className="h-3 w-3" /> Add New Task
      </Button>
      {isEditorOpen && (
        <TaskEditor
          onCancel={handleEditorClose}
          onSubmit={(data) => addTask(data, { onSuccess: handleEditorClose })}
          isLoading={isLoading}
          submitTitle="Add task"
          submitLoadingTitle="Adding new task, please wait..."
        />
      )}
    </li>
  );
}
