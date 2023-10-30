import React from 'react';

import { Header } from '~/common/components';

interface FormLayoutProps {
  children?: React.ReactNode;
  form: React.ReactNode;
  title: string;
  subtitle: string;
}

export function FormLayout({
  children,
  form,
  title,
  subtitle,
}: FormLayoutProps) {
  return (
    <div className="rounded-md border p-5 md:p-8">
      <Header title={title} subtitle={subtitle} />
      {form}
      <div className="pt-6">{children}</div>
    </div>
  );
}
