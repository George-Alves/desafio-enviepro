import type { ComponentProps, ReactNode } from "react";
import { Esqueleto } from "./Esqueleto";
import { GlassCard } from "./GlassCard";

type Props = {
  rotulo: string;
  valor: ReactNode;
  detalhe?: string;
  carregando?: boolean;
  tom?: ComponentProps<typeof GlassCard>["tom"];
  className?: string;
};

export function StatCard({
  rotulo,
  valor,
  detalhe,
  carregando = false,
  tom = "neutro",
  className = "",
}: Props) {
  return (
    <GlassCard tom={tom} eleva className={`flex flex-col gap-1.5 p-6 ${className}`}>
      <p className="text-sm font-medium text-texto-suave">{rotulo}</p>
      {carregando ? (
        <Esqueleto className="mt-1 h-9 w-32" />
      ) : (
        <p className="font-display text-3xl font-semibold tracking-tight">{valor}</p>
      )}
      {detalhe && <p className="text-sm text-texto-suave">{detalhe}</p>}
    </GlassCard>
  );
}
