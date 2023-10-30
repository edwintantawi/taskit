import React, { useState } from 'react';
import { SubmitHandler, useForm } from 'react-hook-form';

import { Button, DatePicker } from '~/common/components';

type FormField = {
  id: string;
  content: string;
  description: string;
  due_date: string | null;
  is_completed: boolean;
};

interface TaskEditorProps {
  onCancel: () => void;
  onSubmit: (data: FormField) => void;
  isLoading?: boolean;
  submitTitle?: string;
  submitLoadingTitle?: string;
  data?: FormField;
}

export function TaskEditor({
  onCancel,
  onSubmit,
  isLoading = false,
  submitTitle = 'Submit',
  submitLoadingTitle = 'Submitting...',
  data,
}: TaskEditorProps) {
  const [dueDate, setDueDate] = useState<string | null>(data?.due_date || null);
  const { register, handleSubmit } = useForm<FormField>({
    defaultValues: data,
  });

  const handleDueDate = (date: string | null) => setDueDate(date);

  const onSubmitForm: SubmitHandler<FormField> = (data) => {
    const { id, content, description, is_completed } = data;
    const newDueDate = dueDate && new Date(dueDate).toISOString();
    onSubmit({
      id,
      content,
      description,
      due_date: newDueDate,
      is_completed,
    });
  };

  return (
    <form className="space-y-3" onSubmit={handleSubmit(onSubmitForm)}>
      <div className="rounded-lg border border-gray-900 p-4">
        <input
          required
          type="text"
          placeholder="Task content"
          className="mb-1 w-full text-sm outline-none"
          {...register('content')}
        />
        <input
          type="text"
          placeholder="Description"
          className="w-full text-xs text-gray-600 outline-none"
          {...register('description')}
        />
        <div className="mt-4">
          <DatePicker onChange={handleDueDate} value={dueDate}>
            Due Date
          </DatePicker>
        </div>
      </div>
      <div className="flex justify-end gap-2">
        <Button
          type="reset"
          variants="outlined"
          size="small"
          onClick={onCancel}
          disabled={isLoading}
        >
          Cancel
        </Button>
        <Button
          variants="contained"
          size="small"
          type="submit"
          disabled={isLoading}
        >
          {isLoading ? submitLoadingTitle : submitTitle}
        </Button>
      </div>
    </form>
  );
}
