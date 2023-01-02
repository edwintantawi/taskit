import React from 'react';

import { Button } from '~/common/components';
import { PlusIcon } from '~/common/components/icons';
import { TaskEditor } from '~/task/components/';

interface TaskListProps {
  children?: React.ReactNode;
}

export function TaskList({ children }: TaskListProps) {
  return (
    <ul>
      {children}
      <li className="mt-4 space-y-4">
        <Button
          variants="normal"
          size="small"
          className="flex items-center gap-1 font-semibold"
        >
          <PlusIcon className="w-3" /> Add New Task
        </Button>
        <TaskEditor />
      </li>
    </ul>
  );
}
