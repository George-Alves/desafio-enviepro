import type { ComponentType, SVGProps } from "react";
import { IconeConfiguracoes, IconeCrm, IconeInicio } from "./icons";

export type ItemMenu = {
  href: string;
  rotulo: string;
  Icone: ComponentType<SVGProps<SVGSVGElement>>;
};

export const ITENS_MENU: ItemMenu[] = [
  { href: "/", rotulo: "Início", Icone: IconeInicio },
  { href: "/crm", rotulo: "CRM", Icone: IconeCrm },
  { href: "/configuracoes", rotulo: "Configurações", Icone: IconeConfiguracoes },
];

/** "/" só é ativo na própria página; as outras rotas também valem para subpáginas. */
export function itemAtivo(href: string, caminho: string): boolean {
  return href === "/" ? caminho === "/" : caminho === href || caminho.startsWith(`${href}/`);
}
