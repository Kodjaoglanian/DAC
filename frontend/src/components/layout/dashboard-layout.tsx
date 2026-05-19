import { TopNav } from './top-nav';
import { Footer } from './footer';

export function DashboardLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen flex flex-col bg-canvas">
      <TopNav />
      <main className="flex-1 max-w-[1200px] w-full mx-auto px-4 lg:px-8 py-8">
        {children}
      </main>
      <Footer />
    </div>
  );
}
