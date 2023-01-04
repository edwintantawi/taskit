import React from 'react';

import { Logo } from './logo';

export function LoadingScreen() {
  return (
    <div className="absolute inset-0 grid cursor-wait place-items-center bg-white p-8">
      <div className="flex flex-col items-center gap-4 text-center">
        <Logo size={44} />
        <h1 className="text-sm leading-6 text-gray-500">
          Loading your awesome <strong>taskit</strong> enviroment,
          <span className="block">please wait...</span>
        </h1>
      </div>
    </div>
  );
}
