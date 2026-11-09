export type StatusCliente = "ativo" | "em_acompanhamento" | "inativo";

export type Cliente = {
  id: number;
  nome: string;
  empresa: string;
  email: string;
  telefone: string;
  status: StatusCliente;
  responsavel: string;
  /** Data no formato AAAA-MM-DD. */
  ultimoContato: string;
  diasSemContato: number;
  /** A API marca true quando passa de 15 dias sem contato. */
  alerta: boolean;
};

export type EtapaFunil = {
  etapa: string;
  quantidade: number;
  valor: number;
};

export type ResumoVendas = {
  total: number;
  quantidade: number;
  ticketMedio: number;
};
