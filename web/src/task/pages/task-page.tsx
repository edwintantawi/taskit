import React from 'react';

import { TaskSection } from '~/task/containers';

export function TaskPage() {
  return (
    <div className="mb-16 space-y-8">
      <TaskSection title="All task" />
    </div>
  );
}
