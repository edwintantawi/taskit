import { AxiosError } from 'axios';
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
    AxiosError<HTTPResponseError>,
    UpdateTaskPayload
  >((payload) => TaskService.updateTask({ ...payload, id: taskId }), {
    onSuccess,
  });

  return mutation;
}
