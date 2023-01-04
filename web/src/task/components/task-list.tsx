import React from 'react';

interface TaskListProps {
  children?: React.ReactNode;
}

export function TaskList({ children }: TaskListProps) {
  return <ul>{children}</ul>;
}
