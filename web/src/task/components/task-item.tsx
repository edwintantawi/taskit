import React from 'react';

import { TaskCheck } from '~/task/components';
import { CalendarIcon } from '~/common/components/icons';

export function TaskItem() {
  return (
    <li className="flex gap-3 border-b py-4">
      <TaskCheck />
      <div className="flex-1 space-y-1">
        <h3 className="text-sm font-semibold">Task content</h3>
        <p className="text-xs text-gray-600">Task description</p>
        <div className="pt-2">
          <p className="flex items-center gap-1 text-xs text-gray-500">
            <CalendarIcon className="h-4" /> 1 Jan 2023 13:00
          </p>
        </div>
      </div>
    </li>
  );
}
