'use client';

import { usePathname } from 'next/navigation';

interface SideMenuItemProps {
  href: string;
  label: string;
}

export function SideMenuItem({ href, label }: SideMenuItemProps) {
  const pathname = usePathname();
  const isActive = pathname === href;

  return (
    <a
      href={href}
      className={`block px-4 py-2 rounded-md transition-colors ${
        isActive ? 'bg-blue-500 text-white' : 'text-gray-700 hover:bg-blue-50 hover:text-blue-600'
      }`}
    >
      {label}
    </a>
  );
}
