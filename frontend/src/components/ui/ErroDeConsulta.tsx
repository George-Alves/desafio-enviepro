import { GlassCard } from "./GlassCard";

type Props = {
  titulo: string;
  erro: unknown;
  onTentarDeNovo: () => void;
  className?: string;
};

/** Diz o que falhou, mostra o motivo e oferece tentar de novo. */
export function ErroDeConsulta({ titulo, erro, onTentarDeNovo, className = "" }: Props) {
  const motivo = erro instanceof Error ? erro.message : "Erro desconhecido.";
  return (
    <GlassCard role="alert" className={`border-[#ff8a7a]/40 p-5 ${className}`}>
      <p className="font-display text-lg font-semibold text-[#ffb4a8]">{titulo}</p>
      <p className="mt-1 text-texto-suave">{motivo}</p>
      <button
        type="button"
        onClick={onTentarDeNovo}
        className="mt-4 rounded-xl bg-white/12 px-4 py-2 font-medium transition-colors hover:bg-white/20 motion-reduce:transition-none"
      >
        Tentar de novo
      </button>
    </GlassCard>
  );
}
