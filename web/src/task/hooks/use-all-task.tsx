import { useQuery } from 'react-query';

import { TaskService } from '~/task/services';

export function useAllTask() {
  const mutation = useQuery('all-task', TaskService.getAllTask);
  return mutation;
}
