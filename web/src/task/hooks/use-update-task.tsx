import { useMutation } from 'react-query';

import { HTTPResponseError, queryClient } from '~/common/libs';
import {
  TaskService,
  UpdateTaskPayload,
  UpdateTaskResponse,
} from '~/task/services';

export function useUpdateTask(taskId: string) {
  const onSuccess = () => {
    queryClient.invalidateQueries({ queryKey: ['all-task'] });
  };

  const mutation = useMutation<
    UpdateTaskResponse,
    HTTPResponseError,
    UpdateTaskPayload
  >((payload) => TaskService.updateTask({ ...payload, id: taskId }), {
    onSuccess,
    onError: (error) => alert(error.error),
  });

  return mutation;
}
