"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { GlassCard } from "@/components/ui/GlassCard";
import { IconeFechar, IconeSetaDireita, IconeSetaEsquerda } from "./icons";
import { ITENS_MENU, itemAtivo } from "./itens-menu";

type Props = {
  id: string;
  /** Só vale no desktop: mostra apenas os ícones. */
  recolhida: boolean;
  /** Só vale no celular: a gaveta está aberta por cima da tela. */
  aberta: boolean;
  onFechar: () => void;
  onAlternarRecolhida: () => void;
};

export function Sidebar({ id, recolhida, aberta, onFechar, onAlternarRecolhida }: Props) {
  const caminho = usePathname();

  // No celular a gaveta fica fora da tela (e invisível, para o Tab não entrar nela) até abrir.
  const posicaoMobile = aberta ? "translate-x-0 visible" : "-translate-x-full invisible";
  const largura = recolhida ? "md:w-[5.5rem]" : "md:w-[16.5rem]";
  // O rótulo some da tela mas continua lido por leitores de tela quando recolhida.
  const rotulo = recolhida ? "md:sr-only" : "";
  const nomeDoBotao = recolhida ? "Expandir menu" : "Recolher menu";

  return (
    <aside
      id={id}
      aria-label="Menu principal"
      className={`fixed inset-y-0 left-0 z-40 w-[17rem] p-3 transition-[transform,visibility,width] duration-300 motion-reduce:transition-none md:visible md:sticky md:top-0 md:z-auto md:h-dvh md:translate-x-0 md:shrink-0 ${posicaoMobile} ${largura}`}
    >
      <GlassCard className="flex h-full flex-col p-3">
        {/* Topo: logo e o botão de recolher, sempre visíveis. Recolhido, o botão fica embaixo do logo. */}
        <div className={`flex items-center gap-3 px-1.5 py-1.5 ${recolhida ? "md:flex-col md:gap-3 md:px-0" : ""}`}>
          <span
            aria-hidden="true"
            className="grid size-10 shrink-0 place-items-center rounded-2xl bg-gradient-to-br from-lavanda to-violeta font-display text-lg font-bold text-white shadow-[0_6px_18px_rgb(139_92_246/0.4)]"
          >
            E
          </span>
          <span className={`whitespace-nowrap font-display text-lg font-semibold leading-tight ${rotulo}`}>
            EnviePro <span className="font-normal text-texto-suave">CRM</span>
          </span>

          <button
            type="button"
            onClick={onAlternarRecolhida}
            aria-label={nomeDoBotao}
            title={nomeDoBotao}
            aria-expanded={!recolhida}
            aria-controls={id}
            className={`hidden size-9 shrink-0 place-items-center rounded-full bg-white/10 text-texto transition-colors hover:bg-white/20 motion-reduce:transition-none md:grid ${recolhida ? "" : "ml-auto"}`}
          >
            {recolhida ? <IconeSetaDireita /> : <IconeSetaEsquerda />}
          </button>

          <button
            type="button"
            onClick={onFechar}
            aria-label="Fechar menu"
            className="ml-auto grid size-9 place-items-center rounded-full bg-white/10 text-texto hover:bg-white/20 md:hidden"
          >
            <IconeFechar />
          </button>
        </div>

        <nav className="mt-6 flex flex-col gap-1.5">
          {ITENS_MENU.map(({ href, rotulo: nome, Icone }) => {
            const ativo = itemAtivo(href, caminho);
            return (
              <Link
                key={href}
                href={href}
                onClick={onFechar}
                aria-current={ativo ? "page" : undefined}
                title={recolhida ? nome : undefined}
                className={`flex items-center gap-3 rounded-full px-4 py-3 text-[0.95rem] font-medium transition-colors motion-reduce:transition-none ${
                  recolhida ? "md:justify-center md:px-0" : ""
                } ${
                  ativo
                    ? "bg-violeta text-white shadow-[0_8px_22px_rgb(139_92_246/0.35)] ring-1 ring-white/25"
                    : "text-texto-suave hover:bg-white/10 hover:text-texto"
                }`}
              >
                <Icone className="shrink-0" />
                <span className={rotulo}>{nome}</span>
              </Link>
            );
          })}
        </nav>
      </GlassCard>
    </aside>
  );
}
