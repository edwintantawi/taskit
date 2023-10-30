import { AxiosError } from 'axios';
import { useMutation } from 'react-query';

import { HTTPResponseError, queryClient } from '~/common/libs';
import { AddTaskPayload, AddTaskResponse, TaskService } from '~/task/services';

export function useAddTask() {
  const onSuccess = (_: AddTaskResponse) => {
    queryClient.invalidateQueries({ queryKey: ['all-task'] });
  };

  const mutation = useMutation<
    AddTaskResponse,
    AxiosError<HTTPResponseError>,
    AddTaskPayload
  >(TaskService.addTask, { onSuccess });

  return mutation;
}
