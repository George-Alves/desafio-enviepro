import { useQuery } from "@tanstack/react-query";
import { listarClientes, obterFunil, obterResumoVendas } from "../services/crm.service";

// Chaves de cache: cada pedido ao servidor tem a sua.
export const chavesCrm = {
  clientes: ["crm", "clientes"] as const,
  funil: ["crm", "funil"] as const,
  resumoVendas: ["crm", "vendas-resumo"] as const,
};

export const useClientes = () => useQuery({ queryKey: chavesCrm.clientes, queryFn: listarClientes });
export const useFunil = () => useQuery({ queryKey: chavesCrm.funil, queryFn: obterFunil });
export const useResumoVendas = () =>
  useQuery({ queryKey: chavesCrm.resumoVendas, queryFn: obterResumoVendas });
