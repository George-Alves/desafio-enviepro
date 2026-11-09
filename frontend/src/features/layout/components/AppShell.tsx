"use client";

import { useEffect, useState, type ReactNode } from "react";
import { GlassCard } from "@/components/ui/GlassCard";
import { IconeMenu } from "./icons";
import { Sidebar } from "./Sidebar";

const ID_MENU = "menu-principal";

/** Moldura de todas as telas: menu lateral + conteúdo. Guarda o estado do menu. */
export function AppShell({ children }: { children: ReactNode }) {
  const [recolhida, setRecolhida] = useState(false);
  const [aberta, setAberta] = useState(false);

  // Gaveta aberta no celular: Esc fecha e a página de trás não rola.
  useEffect(() => {
    if (!aberta) return;
    const aoTeclar = (e: KeyboardEvent) => {
      if (e.key === "Escape") setAberta(false);
    };
    document.addEventListener("keydown", aoTeclar);
    document.body.style.overflow = "hidden";
    return () => {
      document.removeEventListener("keydown", aoTeclar);
      document.body.style.overflow = "";
    };
  }, [aberta]);

  return (
    <div className="flex min-h-dvh">
      <Sidebar
        id={ID_MENU}
        recolhida={recolhida}
        aberta={aberta}
        onFechar={() => setAberta(false)}
        onAlternarRecolhida={() => setRecolhida((v) => !v)}
      />

      {aberta && (
        <div
          aria-hidden="true"
          onClick={() => setAberta(false)}
          className="fixed inset-0 z-30 bg-tinta/60 backdrop-blur-sm md:hidden"
        />
      )}

      <div className="flex min-w-0 flex-1 flex-col">
        <GlassCard className="sticky top-3 z-20 mx-3 mt-3 flex items-center gap-3 px-3 py-2 md:hidden">
          <button
            type="button"
            onClick={() => setAberta(true)}
            aria-label="Abrir menu"
            aria-expanded={aberta}
            aria-controls={ID_MENU}
            className="grid size-10 place-items-center rounded-lg hover:bg-white/10"
          >
            <IconeMenu />
          </button>
          <span className="font-display text-lg font-semibold">
            EnviePro <span className="font-normal text-texto-suave">CRM</span>
          </span>
        </GlassCard>

        <main className="flex-1 px-5 py-6 md:px-10 md:py-10">{children}</main>
      </div>
    </div>
  );
}
