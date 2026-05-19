# ProjectTracker Frontend

Frontend da Plataforma de Registro e Acompanhamento de Projetos.

## Stack

- Next.js 14 (App Router)
- React 18 + TypeScript
- Tailwind CSS 3
- Radix UI Primitives
- Axios + React Hook Form + Zod
- Recharts (gráficos)
- next-themes (tema claro/escuro)
- sonner (toasts)

## Setup Local

```bash
cd frontend
npm install
cp .env.example .env.local
npm run dev
```

## Build

```bash
npm run build
```

## Docker

```bash
# Build e run
docker compose up -d --build

# Stop
docker compose down
```

A aplicação será exposta na porta `3000`.

## Variáveis de Ambiente

| Variável | Padrão | Descrição |
|---|---|---|
| `NEXT_PUBLIC_API_URL` | http://localhost:8080/api/v1 | URL base da API |
| `NEXT_PUBLIC_APP_NAME` | ProjectTracker | Nome do app |
| `NEXT_PUBLIC_APP_URL` | http://localhost:3000 | URL do app |

## Estrutura

```
src/
├── app/              # Páginas Next.js
├── components/       # Componentes React
│   ├── layout/       # Layout (TopNav, Footer, DashboardLayout)
│   ├── shared/       # Design System (Button, Input, Badge, Avatar)
│   └── ...           # Domínio (projects, tasks, dashboard, members)
├── lib/
│   ├── api/          # Cliente Axios + funções de API
│   ├── hooks/        # Hooks customizados
│   └── utils.ts      # Utilitários (cn, formatDate)
├── providers/        # Contextos (Auth, Theme, Toast)
├── types/            # Tipos TypeScript
└── middleware.ts     # Proteção de rotas
```

## Segurança

- Headers de segurança (CSP, X-Frame-Options, etc.)
- Middleware de autenticação
- Interceptors Axios (401 → redirect)
- Runtime Docker sem root (`USER node`)
