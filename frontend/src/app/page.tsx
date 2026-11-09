import Link from "next/link";
import { GlassCard } from "@/components/ui/GlassCard";

export default function Inicio() {
  return (
    <div className="max-w-2xl">
      <h1 className="font-display text-4xl font-semibold tracking-tight md:text-5xl">
        Bem-vindo ao EnviePro CRM
      </h1>
      <p className="mt-4 max-w-prose text-lg leading-relaxed text-texto-suave">
        Veja quem precisa de contato, em que etapa está cada negociação e quanto já foi vendido no mês.
      </p>

      <GlassCard className="mt-8 p-6">
        <h2 className="font-display text-xl font-semibold">Painel do CRM</h2>
        <p className="mt-2 text-texto-suave">
          Clientes com acompanhamento, funil de vendas e resumo do mês, todos na mesma tela.
        </p>
        <Link
          href="/crm"
          className="mt-5 inline-flex rounded-xl bg-acafrao px-5 py-2.5 font-semibold text-tinta transition-colors hover:bg-[#f0b657] motion-reduce:transition-none"
        >
          Abrir o CRM
        </Link>
      </GlassCard>
    </div>
  );
}
