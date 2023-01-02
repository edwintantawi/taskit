import React from 'react';

import { TaskItem, TaskList } from '~/task/components';

export function TaskPage() {
  return (
    <div>
      <TaskList>
        <TaskItem />
        <TaskItem />
        <TaskItem />
        <TaskItem />
        <TaskItem />
      </TaskList>
    </div>
  );
}
