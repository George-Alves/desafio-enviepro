import { GlassCard } from "@/components/ui/GlassCard";

export default function ConfiguracoesPage() {
  return (
    <div className="mx-auto max-w-5xl">
      <h1 className="font-display text-4xl font-semibold tracking-tight">Configurações</h1>
      <GlassCard tom="indigo" className="mt-6 max-w-2xl p-7">
        <p className="text-texto-suave">Ainda não há nada para configurar neste sistema.</p>
      </GlassCard>
    </div>
  );
}
