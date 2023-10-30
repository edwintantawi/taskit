import { useMutation } from 'react-query';

import { HTTPResponseError, queryClient } from '~/common/libs';
import { AddTaskPayload, AddTaskResponse, TaskService } from '~/task/services';

export function useAddTask() {
  const onSuccess = (_: AddTaskResponse) => {
    queryClient.invalidateQueries({ queryKey: ['all-task'] });
  };

  const mutation = useMutation<
    AddTaskResponse,
    HTTPResponseError,
    AddTaskPayload
  >(TaskService.addTask, {
    onSuccess,
    onError: (error) => alert(error.error),
  });

  return mutation;
}
