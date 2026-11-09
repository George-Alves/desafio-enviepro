const moeda = new Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" });
const moedaInteira = new Intl.NumberFormat("pt-BR", {
  style: "currency",
  currency: "BRL",
  maximumFractionDigits: 0,
});

export const formatarMoeda = (valor: number) => moeda.format(valor);
export const formatarMoedaInteira = (valor: number) => moedaInteira.format(valor);

/** "2026-11-09" -> "09/11/2026", sem passar por Date (evita erro de fuso horário). */
export function formatarData(iso: string): string {
  const [ano, mes, dia] = iso.split("-");
  return `${dia}/${mes}/${ano}`;
}
