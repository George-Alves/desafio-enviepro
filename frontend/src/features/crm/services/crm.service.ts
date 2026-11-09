import { apiGet } from "@/lib/api/client";
import type { Cliente, EtapaFunil, ResumoVendas } from "../types";

// Formato que a API devolve (snake_case). Fica só aqui: o resto do front usa os tipos de ../types.
type ClienteApi = {
  id: number;
  nome: string;
  empresa: string;
  email: string;
  telefone: string;
  status: Cliente["status"];
  responsavel: string;
  ultimo_contato: string;
  dias_sem_contato: number;
  alerta: boolean;
};

type ResumoApi = { total: number; quantidade: number; ticket_medio: number };

export async function listarClientes(): Promise<Cliente[]> {
  const dados = await apiGet<ClienteApi[]>("/clientes");
  return dados.map((c) => ({
    id: c.id,
    nome: c.nome,
    empresa: c.empresa,
    email: c.email,
    telefone: c.telefone,
    status: c.status,
    responsavel: c.responsavel,
    ultimoContato: c.ultimo_contato,
    diasSemContato: c.dias_sem_contato,
    alerta: c.alerta,
  }));
}

export function obterFunil(): Promise<EtapaFunil[]> {
  return apiGet<EtapaFunil[]>("/funil");
}

export async function obterResumoVendas(): Promise<ResumoVendas> {
  const r = await apiGet<ResumoApi>("/vendas/resumo");
  return { total: r.total, quantidade: r.quantidade, ticketMedio: r.ticket_medio };
}
