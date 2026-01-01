import { type ReactNode } from 'react';

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  return (
    <div className="min-h-screen bg-gray-100 flex">
      {/* Side Menu */}
      <aside className="w-64 bg-white shadow-md">
        <div className="p-6">
          <h1 className="text-2xl font-bold text-gray-900 mb-8">Platform</h1>
          <nav className="space-y-2">
            <a
              href="/me"
              className="block px-4 py-2 text-gray-700 hover:bg-blue-50 hover:text-blue-600 rounded-md transition-colors"
            >
              My Profile
            </a>
            <a
              href="/users"
              className="block px-4 py-2 text-gray-700 hover:bg-blue-50 hover:text-blue-600 rounded-md transition-colors"
            >
              Workspace Users
            </a>
          </nav>
        </div>
        <div className="absolute bottom-0 w-64 p-6 border-t border-gray-200">
          <a
            href="/auth/logout"
            className="block w-full bg-red-500 text-white py-2 px-4 rounded-md hover:bg-red-600 text-center"
          >
            Logout
          </a>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 p-8">{children}</main>
    </div>
  );
}
