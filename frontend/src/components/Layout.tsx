import type { ReactNode } from 'react';
import { SideMenuItem } from './SideMenuItem';

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  return (
    <div className="min-h-screen bg-gray-100 flex">
      {/* Side Menu */}
      <aside className="w-64 bg-gray-50 shadow-lg flex flex-col">
        <div className="px-6 py-6 border-b border-gray-200">
          <h1 className="text-2xl font-bold text-gray-900">Platform</h1>
        </div>
        <nav className="flex-1 p-4 space-y-2">
          <SideMenuItem href="/me" label="My Profile" />
          <SideMenuItem href="/users" label="Workspace Users" />
        </nav>
        <div className="p-4 border-t border-gray-200">
          <a
            href="/auth/logout"
            className="block w-full bg-red-500 text-white py-2 px-4 rounded-md hover:bg-red-600 text-center"
          >
            Logout
          </a>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1">{children}</main>
    </div>
  );
}
