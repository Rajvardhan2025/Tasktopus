import { Outlet, Link, useLocation } from 'react-router-dom';
import { LayoutDashboard, FolderKanban } from 'lucide-react';
import { cn } from '@/lib/utils';

export function Layout() {
  const location = useLocation();

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <Link to="/" className="flex items-center gap-2">
              <LayoutDashboard className="w-6 h-6" />
              <h1 className="text-xl font-bold">Project Management</h1>
            </Link>

            <nav className="flex items-center gap-6">
              <Link
                to="/projects"
                className={cn(
                  "flex items-center gap-2 text-sm font-medium transition-colors hover:text-primary",
                  location.pathname === '/projects' ? 'text-primary' : 'text-muted-foreground'
                )}
              >
                <FolderKanban className="w-4 h-4" />
                Projects
              </Link>
            </nav>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        <Outlet />
      </main>
    </div>
  );
}
