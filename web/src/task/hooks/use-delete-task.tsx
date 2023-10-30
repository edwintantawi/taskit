import { AxiosError, AxiosResponse } from 'axios';
import { useMutation } from 'react-query';

import { HTTPResponseError, queryClient } from '~/common/libs';
import { TaskService } from '~/task/services';

export function useDeleteTask(taskId: string) {
  const onSuccess = () => {
    queryClient.invalidateQueries({ queryKey: ['all-task'] });
  };

  const mutation = useMutation<AxiosResponse, AxiosError<HTTPResponseError>>(
    () => TaskService.deleteTask(taskId),
    { onSuccess }
  );

  return mutation;
}
