import React from 'react';

import { TaskItem, TaskList } from '~/task/components';
import { TaskSectionAction } from '~/task/containers';
import { useAllTask } from '~/task/hooks';

export function TaskSection({ title }: { title: string }) {
  const { data: tasks, isLoading } = useAllTask();

  return (
    <section>
      <h2 className="mb-4 text-xl font-bold">{title}</h2>
      {isLoading && <p>Loading...</p>}
      <TaskList>
        {tasks?.data.payload.map((task) => {
          return <TaskItem key={task.id} {...task} />;
        })}
        <TaskSectionAction />
      </TaskList>
    </section>
  );
}
