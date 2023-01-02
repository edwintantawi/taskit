import React from 'react';

import { Button } from '~/common/components';
import { CalendarIcon } from '~/common/components/icons';

export function TaskEditor() {
  return (
    <div className="space-y-3">
      <div className="rounded-lg border border-gray-900 p-4">
        <input
          required
          type="text"
          placeholder="Task content"
          className="mb-1 w-full text-sm outline-none"
        />
        <input
          type="text"
          placeholder="Description"
          className="w-full text-xs text-gray-600 outline-none"
        />
        <div className="mt-4">
          <Button
            variants="outlined"
            size="small"
            className="flex items-center gap-1 border-gray-500 py-1 text-gray-500"
          >
            <CalendarIcon className="w-3" /> Due Date
          </Button>
        </div>
      </div>
      <div className="flex justify-end gap-2">
        <Button variants="outlined" size="small">
          Cancel
        </Button>
        <Button variants="contained" size="small">
          Add task
        </Button>
      </div>
    </div>
  );
}
