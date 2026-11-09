"use client";

import { ErroDeConsulta } from "@/components/ui/ErroDeConsulta";
import { StatCard } from "@/components/ui/StatCard";
import { formatarMoeda } from "@/lib/formatar";
import { useClientes, useResumoVendas } from "../hooks/useCrm";

export function CardsResumo() {
  const resumo = useResumoVendas();
  const clientes = useClientes();

  if (resumo.isError) {
    return (
      <ErroDeConsulta
        titulo="Não foi possível carregar o resumo de vendas"
        erro={resumo.error}
        onTentarDeNovo={() => resumo.refetch()}
      />
    );
  }

  const ativos = clientes.data?.filter((c) => c.status === "ativo").length;

  return (
    <section aria-label="Resumo de vendas do mês" className="grid grid-cols-2 gap-4 xl:grid-cols-4">
      <StatCard
        rotulo="Total de vendas no mês"
        valor={resumo.data && formatarMoeda(resumo.data.total)}
        carregando={resumo.isPending}
        destaque
        className="col-span-2 sm:col-span-1"
      />
      <StatCard
        rotulo="Ticket médio"
        valor={resumo.data && formatarMoeda(resumo.data.ticketMedio)}
        carregando={resumo.isPending}
        className="col-span-2 sm:col-span-1"
      />
      <StatCard
        rotulo="Vendas no mês"
        valor={resumo.data?.quantidade}
        detalhe={
          resumo.data ? (resumo.data.quantidade === 1 ? "venda fechada" : "vendas fechadas") : undefined
        }
        carregando={resumo.isPending}
      />
      <StatCard
        rotulo="Clientes ativos"
        valor={clientes.isError ? "—" : ativos}
        detalhe={clientes.data ? `de ${clientes.data.length} clientes` : undefined}
        carregando={clientes.isPending}
      />
    </section>
  );
}
