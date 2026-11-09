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
  const largura = recolhida ? "md:w-[5.25rem]" : "md:w-[16.5rem]";
  // O rótulo some da tela mas continua lido por leitores de tela quando recolhida.
  const rotulo = recolhida ? "md:sr-only" : "";

  return (
    <aside
      id={id}
      aria-label="Menu principal"
      className={`fixed inset-y-0 left-0 z-40 w-[17rem] p-3 transition-[transform,visibility,width] duration-300 motion-reduce:transition-none md:visible md:sticky md:top-0 md:z-auto md:h-dvh md:translate-x-0 md:shrink-0 ${posicaoMobile} ${largura}`}
    >
      <GlassCard className="flex h-full flex-col p-3">
        <div className="flex items-center gap-3 px-1.5 py-1.5">
          <span
            aria-hidden="true"
            className="grid size-9 shrink-0 place-items-center rounded-xl bg-acafrao font-display text-lg font-bold text-tinta"
          >
            E
          </span>
          <span className={`whitespace-nowrap font-display text-lg font-semibold leading-tight ${rotulo}`}>
            EnviePro <span className="font-normal text-texto-suave">CRM</span>
          </span>
          <button
            type="button"
            onClick={onFechar}
            aria-label="Fechar menu"
            className="ml-auto grid size-9 place-items-center rounded-lg text-texto-suave hover:bg-white/10 hover:text-texto md:hidden"
          >
            <IconeFechar />
          </button>
        </div>

        <nav className="mt-5 flex flex-col gap-1">
          {ITENS_MENU.map(({ href, rotulo: nome, Icone }) => {
            const ativo = itemAtivo(href, caminho);
            return (
              <Link
                key={href}
                href={href}
                onClick={onFechar}
                aria-current={ativo ? "page" : undefined}
                title={recolhida ? nome : undefined}
                className={`relative flex items-center gap-3 rounded-xl px-3.5 py-2.5 text-[0.95rem] font-medium transition-colors motion-reduce:transition-none ${
                  ativo
                    ? "bg-white/15 text-texto"
                    : "text-texto-suave hover:bg-white/10 hover:text-texto"
                }`}
              >
                {ativo && (
                  <span
                    aria-hidden="true"
                    className="absolute inset-y-2 left-0 w-[3px] rounded-full bg-acafrao"
                  />
                )}
                <Icone className="shrink-0" />
                <span className={rotulo}>{nome}</span>
              </Link>
            );
          })}
        </nav>

        <button
          type="button"
          onClick={onAlternarRecolhida}
          aria-label={recolhida ? "Expandir menu" : "Recolher menu"}
          aria-expanded={!recolhida}
          aria-controls={id}
          className="mt-auto hidden items-center gap-3 rounded-xl px-3.5 py-2.5 text-[0.95rem] font-medium text-texto-suave transition-colors hover:bg-white/10 hover:text-texto md:flex"
        >
          {recolhida ? (
            <IconeSetaDireita className="shrink-0" />
          ) : (
            <IconeSetaEsquerda className="shrink-0" />
          )}
          <span className={rotulo}>{recolhida ? "Expandir menu" : "Recolher menu"}</span>
        </button>
      </GlassCard>
    </aside>
  );
}
