"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState, type ReactNode } from "react";

export function ReactQueryProvider({ children }: { children: ReactNode }) {
  // useState garante um cliente por navegador, e não um novo a cada renderização.
  const [cliente] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 30_000, // dentro de 30s reaproveita o que já buscou
            retry: 1,
            refetchOnWindowFocus: true, // voltar para a aba atualiza os números
          },
        },
      }),
  );
  return <QueryClientProvider client={cliente}>{children}</QueryClientProvider>;
}
