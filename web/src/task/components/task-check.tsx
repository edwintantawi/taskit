import React, { useState } from 'react';

import { CheckIcon } from '~/common/components/icons';

export function TaskCheck() {
  const [isChecked, setIsChecked] = useState(false);
  const handleChange = () => setIsChecked((prev) => !prev);

  return (
    <div className="relative mt-1">
      <input
        type="checkbox"
        className="peer absolute top-0 z-50 aspect-square h-5 w-5 cursor-pointer opacity-0"
        checked={isChecked}
        onChange={handleChange}
      />
      <div
        className={`grid h-5 w-5 place-items-center rounded-full border border-gray-400 text-white ${
          isChecked ? 'bg-gray-900' : 'bg-white peer-hover:text-gray-400'
        }`}
      >
        <CheckIcon className="h-3 w-3" />
      </div>
    </div>
  );
}
