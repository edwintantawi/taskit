import React, { useRef, useState } from 'react';
import { CalendarIcon } from '@heroicons/react/24/outline';
import { format } from 'date-fns';

import { Button } from '~/common/components';

interface DatePickerProps {
  children?: React.ReactNode;
  onChange: (date: string | null) => void;
  value?: string | null;
}

export function DatePicker({ children, onChange, value }: DatePickerProps) {
  const [dueDate, setDueDate] = useState<string | null>(value ?? null);
  const dateTimeRef = useRef<HTMLInputElement | null>(null);

  const handleOpenDateTimePicker = () => dateTimeRef.current?.focus();

  return (
    <Button
      type="button"
      variants="outlined"
      size="small"
      className="relative flex items-center gap-1 border-gray-500 py-1 text-gray-500"
      onClick={handleOpenDateTimePicker}
    >
      <CalendarIcon className="h-3 w-3" />
      {dueDate ? format(new Date(dueDate), 'dd MMMM yyyy') : children}
      <input
        type="date"
        onChange={(event) => {
          const selectedDate = event.target.value || null;
          onChange(selectedDate);
          setDueDate(selectedDate);
        }}
        onFocus={(e) => e.target.showPicker()}
        className="absolute left-0 h-0 w-0"
        ref={dateTimeRef}
      />
    </Button>
  );
}
