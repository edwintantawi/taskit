import { AxiosResponse } from 'axios';

import { http, HTTPResponseSuccess } from '~/common/libs';
import { Task } from '~/task/entity';

export type AddTaskPayload = {
  content: string;
  description: string;
  due_date: string | null;
};

export type UpdateTaskPayload = {
  id: string;
  content: string;
  description: string;
  is_completed: boolean;
  due_date: string | null;
};

export type GetAllTaskResponse = AxiosResponse<HTTPResponseSuccess<Task[]>>;
export type AddTaskResponse = AxiosResponse<
  HTTPResponseSuccess<{ id: string }>
>;
export type UpdateTaskResponse = AxiosResponse<
  HTTPResponseSuccess<{ id: string }>
>;

export class TaskService {
  static getAllTask(): Promise<GetAllTaskResponse> {
    return http('/tasks', {
      method: 'GET',
    });
  }

  static addTask(task: AddTaskPayload): Promise<AddTaskResponse> {
    return http('/tasks', {
      method: 'POST',
      data: {
        content: task.content,
        description: task.description,
        due_date: task.due_date,
      },
    });
  }

  static updateTask(task: UpdateTaskPayload): Promise<UpdateTaskResponse> {
    return http(`/tasks/${task.id}`, {
      method: 'PUT',
      data: {
        content: task.content,
        description: task.description,
        due_date: task.due_date,
        is_completed: task.is_completed,
      },
    });
  }

  static deleteTask(id: string) {
    return http(`/tasks/${id}`, {
      method: 'DELETE',
    });
  }
}
