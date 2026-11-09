import type { ComponentProps } from "react";

/** Bloco cinza pulsante que ocupa o lugar do conteúdo enquanto ele carrega. */
export function Esqueleto({ className = "", ...props }: ComponentProps<"div">) {
  return (
    <div
      aria-hidden="true"
      className={`animate-pulse rounded-lg bg-white/10 motion-reduce:animate-none ${className}`}
      {...props}
    />
  );
}
