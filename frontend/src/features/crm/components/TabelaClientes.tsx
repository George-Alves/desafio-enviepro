"use client";

import { useMemo, useState } from "react";
import { ErroDeConsulta } from "@/components/ui/ErroDeConsulta";
import { Esqueleto } from "@/components/ui/Esqueleto";
import { GlassCard } from "@/components/ui/GlassCard";
import { formatarData } from "@/lib/formatar";
import { useClientes } from "../hooks/useCrm";
import type { StatusCliente } from "../types";

const STATUS: Record<StatusCliente, { rotulo: string; cor: string }> = {
  ativo: { rotulo: "Ativo", cor: "text-[#7fe0d0]" },
  em_acompanhamento: { rotulo: "Em acompanhamento", cor: "text-[#f2c26b]" },
  inativo: { rotulo: "Inativo", cor: "text-texto-suave" },
};

/** Tira acentos e maiúsculas: "JOSÉ" encontra "jose". */
const normalizar = (s: string) =>
  s.normalize("NFD").replace(/\p{Diacritic}/gu, "").toLowerCase();

function textoDias(dias: number): string {
  if (dias <= 0) return "hoje";
  if (dias === 1) return "ontem";
  return `há ${dias} dias`;
}

const CAMPO =
  "rounded-xl border border-white/15 bg-white/10 px-3.5 py-2.5 text-texto placeholder:text-texto-suave";

export function TabelaClientes() {
  const clientes = useClientes();
  const [busca, setBusca] = useState("");
  const [status, setStatus] = useState<"todos" | StatusCliente>("todos");

  const filtrados = useMemo(() => {
    const termo = normalizar(busca.trim());
    return (clientes.data ?? []).filter(
      (c) => (status === "todos" || c.status === status) && normalizar(c.nome).includes(termo),
    );
  }, [clientes.data, busca, status]);

  if (clientes.isError) {
    return (
      <ErroDeConsulta
        titulo="Não foi possível carregar os clientes"
        erro={clientes.error}
        onTentarDeNovo={() => clientes.refetch()}
      />
    );
  }

  const limparFiltros = () => {
    setBusca("");
    setStatus("todos");
  };

  return (
    <GlassCard className="p-5 md:p-6" role="region" aria-labelledby="titulo-clientes">
      <h2 id="titulo-clientes" className="font-display text-xl font-semibold">
        Clientes
      </h2>
      <p className="mt-1 text-sm text-texto-suave">
        Clientes sem contato há mais de 15 dias aparecem em vermelho.
      </p>

      <div className="mt-5 flex flex-col gap-3 sm:flex-row">
        <label className="flex flex-1 flex-col gap-1.5 text-sm font-medium text-texto-suave">
          Buscar por nome
          <input
            type="search"
            value={busca}
            onChange={(e) => setBusca(e.target.value)}
            placeholder="Ex.: Ana"
            className={`${CAMPO} text-base font-normal`}
          />
        </label>
        <label className="flex flex-col gap-1.5 text-sm font-medium text-texto-suave sm:w-56">
          Status
          <select
            value={status}
            onChange={(e) => setStatus(e.target.value as "todos" | StatusCliente)}
            className={`${CAMPO} text-base font-normal [color-scheme:dark]`}
          >
            <option value="todos">Todos</option>
            {Object.entries(STATUS).map(([valor, { rotulo }]) => (
              <option key={valor} value={valor}>
                {rotulo}
              </option>
            ))}
          </select>
        </label>
      </div>

      <p aria-live="polite" className="mt-4 text-sm text-texto-suave">
        {clientes.isSuccess &&
          `${filtrados.length} ${filtrados.length === 1 ? "cliente" : "clientes"}`}
      </p>

      {clientes.isPending ? (
        <div className="mt-3 flex flex-col gap-2">
          {Array.from({ length: 6 }, (_, i) => (
            <Esqueleto key={i} className="h-12 w-full" />
          ))}
        </div>
      ) : filtrados.length === 0 ? (
        <div className="mt-3 rounded-xl bg-white/5 p-5 text-texto-suave">
          <p>Nenhum cliente encontrado com esses filtros.</p>
          <button
            type="button"
            onClick={limparFiltros}
            className="mt-3 rounded-xl bg-white/12 px-4 py-2 font-medium text-texto transition-colors hover:bg-white/20 motion-reduce:transition-none"
          >
            Limpar filtros
          </button>
        </div>
      ) : (
        <div className="mt-3 overflow-x-auto">
          <table className="w-full min-w-[34rem] border-separate border-spacing-y-1.5 text-left">
            <caption className="sr-only">Lista de clientes com status, responsável e último contato</caption>
            <thead>
              <tr className="text-sm text-texto-suave">
                <th scope="col" className="px-4 pb-1 font-medium">Cliente</th>
                <th scope="col" className="px-4 pb-1 font-medium">Status</th>
                <th scope="col" className="px-4 pb-1 font-medium">Responsável</th>
                <th scope="col" className="px-4 pb-1 font-medium">Último contato</th>
              </tr>
            </thead>
            <tbody>
              {filtrados.map((c) => {
                const linha = c.alerta ? "bg-[#ff6b5b]/15" : "bg-white/5";
                return (
                  <tr key={c.id} className={linha}>
                    <th
                      scope="row"
                      className={`rounded-l-xl px-4 py-3 font-medium ${
                        c.alerta ? "shadow-[inset_3px_0_0_#ff8a7a]" : ""
                      }`}
                    >
                      {c.nome}
                      <span className="block text-sm font-normal text-texto-suave">{c.empresa}</span>
                    </th>
                    <td className={`px-4 py-3 ${STATUS[c.status].cor}`}>{STATUS[c.status].rotulo}</td>
                    <td className="px-4 py-3">{c.responsavel}</td>
                    <td className="rounded-r-xl px-4 py-3">
                      <span className={c.alerta ? "font-semibold text-[#ffb4a8]" : ""}>
                        {formatarData(c.ultimoContato)}
                      </span>
                      <span className={`block text-sm ${c.alerta ? "text-[#ffb4a8]" : "text-texto-suave"}`}>
                        {textoDias(c.diasSemContato)}
                        {c.alerta && <span className="sr-only"> (atenção: sem contato há mais de 15 dias)</span>}
                      </span>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </GlassCard>
  );
}
