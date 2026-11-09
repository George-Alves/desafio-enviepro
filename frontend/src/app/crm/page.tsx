import { PainelCrm } from "@/features/crm";

export default function CrmPage() {
  return (
    <div className="mx-auto max-w-5xl">
      <h1 className="font-display text-4xl font-semibold tracking-tight">CRM</h1>
      <p className="mt-2 mb-6 text-texto-suave">Vendas do mês, funil e acompanhamento dos clientes.</p>
      <PainelCrm />
    </div>
  );
}
