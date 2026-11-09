"use client";

import { ErroDeConsulta } from "@/components/ui/ErroDeConsulta";
import { Esqueleto } from "@/components/ui/Esqueleto";
import { GlassCard } from "@/components/ui/GlassCard";
import { formatarMoedaInteira } from "@/lib/formatar";
import { useFunil } from "../hooks/useCrm";

export function FunilVendas() {
  const funil = useFunil();

  if (funil.isError) {
    return (
      <ErroDeConsulta
        titulo="Não foi possível carregar o funil de vendas"
        erro={funil.error}
        onTentarDeNovo={() => funil.refetch()}
      />
    );
  }

  const etapas = funil.data ?? [];
  const maiorValor = Math.max(0, ...etapas.map((e) => e.valor));
  const semOportunidades = funil.isSuccess && etapas.every((e) => e.quantidade === 0);

  return (
    <GlassCard className="p-6 md:p-7" aria-labelledby="titulo-funil" role="region">
      <h2 id="titulo-funil" className="font-display text-xl font-semibold">
        Funil de vendas
      </h2>
      <p className="mt-1 text-sm text-texto-suave">
        Oportunidades por etapa. O tamanho da barra mostra o valor em negociação.
      </p>

      {semOportunidades ? (
        <p className="mt-5 text-texto-suave">Ainda não há oportunidades no funil.</p>
      ) : (
        <ul className="mt-5 flex flex-col gap-4">
          {funil.isPending
            ? Array.from({ length: 5 }, (_, i) => (
                <li key={i}>
                  <Esqueleto className="h-10 w-full" />
                </li>
              ))
            : etapas.map((e) => {
                const largura = maiorValor > 0 ? Math.max(2, (e.valor / maiorValor) * 100) : 0;
                const fechada = e.etapa === "Fechado";
                return (
                  <li key={e.etapa}>
                    <div className="flex flex-wrap items-baseline justify-between gap-x-4">
                      <span className="font-medium">{e.etapa}</span>
                      <span className="text-sm text-texto-suave">
                        <span className="font-semibold text-texto">
                          {e.quantidade} {e.quantidade === 1 ? "oportunidade" : "oportunidades"}
                        </span>
                        {" · "}
                        {formatarMoedaInteira(e.valor)}
                      </span>
                    </div>
                    <div className="mt-1.5 h-3 overflow-hidden rounded-full bg-white/10" aria-hidden="true">
                      <div
                        className={`h-full rounded-full bg-gradient-to-r ${
                          fechada ? "from-lavanda to-nevoa" : "from-violeta to-lavanda"
                        }`}
                        style={{ width: `${largura}%` }}
                      />
                    </div>
                  </li>
                );
              })}
        </ul>
      )}
    </GlassCard>
  );
}
