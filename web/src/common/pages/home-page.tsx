import React from 'react';

import { Button, Footer, Header } from '~/common/components';

export function HomePage() {
  return (
    <div className="flex flex-1 flex-col justify-between">
      <div className="py-14 text-center">
        <Header
          title="Organize Your Task and Life More Easier Than You Think."
          subtitle="Become focused, organize and enjoy your task with Taskit."
        />
        <Button
          variants="contained"
          size="large"
          className="shadow-xl shadow-gray-300 transition-shadow duration-500 hover:shadow-none"
        >
          Start Your Journey Here
        </Button>
      </div>
      <img src="/figure.svg" alt="" className="aspect-auto w-full" />
      <Footer />
    </div>
  );
}
