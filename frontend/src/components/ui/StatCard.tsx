import type { ReactNode } from "react";
import { Esqueleto } from "./Esqueleto";
import { GlassCard } from "./GlassCard";

type Props = {
  rotulo: string;
  valor: ReactNode;
  detalhe?: string;
  carregando?: boolean;
  /** Destaca o número na cor de acento. Use em um card só. */
  destaque?: boolean;
  className?: string;
};

export function StatCard({
  rotulo,
  valor,
  detalhe,
  carregando = false,
  destaque = false,
  className = "",
}: Props) {
  return (
    <GlassCard className={`flex flex-col gap-1.5 p-5 ${className}`}>
      <p className="text-sm font-medium text-texto-suave">{rotulo}</p>
      {carregando ? (
        <Esqueleto className="mt-1 h-9 w-32" />
      ) : (
        <p
          className={`font-display text-3xl font-semibold tracking-tight ${destaque ? "text-acafrao" : ""}`}
        >
          {valor}
        </p>
      )}
      {detalhe && <p className="text-sm text-texto-suave">{detalhe}</p>}
    </GlassCard>
  );
}
