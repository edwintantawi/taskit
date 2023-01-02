import React from 'react';

import { TaskCheck } from '~/task/components';
import { CalendarIcon, DeleteIcon, EditIcon } from '~/common/components/icons';

export function TaskItem() {
  return (
    <li className="group flex gap-3 border-b py-4">
      <TaskCheck />
      <div className="relative w-[calc(100%-32px)] flex-1">
        <h3 className="mb-1 w-full break-words text-sm leading-[1.4rem]">
          Task content
        </h3>
        <p className="mb-2 text-xs text-gray-600">Task description</p>
        <p className="flex items-center gap-1 text-xs text-gray-500">
          <CalendarIcon className="h-4" /> 1 Jan 2023 13:00
        </p>

        <div className=" absolute top-0 right-0 flex h-10 text-sm text-gray-500 group-hover:visible">
          <div className="w-10  bg-gradient-to-l from-white to-transparent" />
          <div className="flex items-center gap-1 bg-white pl-3">
            <button className="rounded-md p-1 hover:bg-gray-100">
              <EditIcon className="h-5 w-5" />
            </button>
            <button className="rounded-md p-1 text-red-500 hover:bg-gray-100">
              <DeleteIcon className="h-5 w-5" />
            </button>
          </div>
        </div>
      </div>
    </li>
  );
}
