import { CardsResumo } from "./CardsResumo";
import { FunilVendas } from "./FunilVendas";
import { TabelaClientes } from "./TabelaClientes";

export function PainelCrm() {
  return (
    <div className="flex flex-col gap-6">
      <CardsResumo />
      <FunilVendas />
      <TabelaClientes />
    </div>
  );
}
