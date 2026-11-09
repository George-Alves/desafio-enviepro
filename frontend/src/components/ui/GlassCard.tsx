import type { ComponentProps } from "react";

/** Painel de vidro (glassmorphism). Toda tela usa este componente em vez de repetir o estilo. */
export function GlassCard({ className = "", ...props }: ComponentProps<"div">) {
  return <div className={`vidro rounded-2xl ${className}`} {...props} />;
}
