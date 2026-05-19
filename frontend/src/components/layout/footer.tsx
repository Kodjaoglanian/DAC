export function Footer() {
  return (
    <footer className="bg-surface-dark text-on-dark-soft font-body text-caption py-6 px-8">
      <div className="max-w-[1200px] mx-auto flex flex-col sm:flex-row items-center justify-between gap-4">
        <div className="flex items-center gap-2">
          <span className="text-on-dark font-medium">ProjectTracker</span>
          <span className="text-muted-soft">Plataforma de Registro e Acompanhamento de Projetos</span>
        </div>
        <p className="text-muted-soft">
          &copy; {new Date().getFullYear()} ProjectTracker
        </p>
      </div>
    </footer>
  );
}
