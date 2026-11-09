/** Endereço da API Go. Pode ser trocado com NEXT_PUBLIC_API_URL (ex.: se a porta 8080 estiver ocupada). */
export const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export class ErroDeApi extends Error {
  constructor(
    mensagem: string,
    readonly status?: number,
  ) {
    super(mensagem);
    this.name = "ErroDeApi";
  }
}

/** Único ponto do front que fala com a API: monta a URL e padroniza os erros. */
export async function apiGet<T>(caminho: string): Promise<T> {
  let resposta: Response;
  try {
    resposta = await fetch(`${API_URL}${caminho}`, { headers: { Accept: "application/json" } });
  } catch {
    throw new ErroDeApi(`Não foi possível falar com a API em ${API_URL}. Ela está rodando?`);
  }
  if (!resposta.ok) {
    throw new ErroDeApi(`A API respondeu com erro ${resposta.status} em ${caminho}.`, resposta.status);
  }
  try {
    return (await resposta.json()) as T;
  } catch {
    throw new ErroDeApi(`A resposta de ${caminho} não é um JSON válido. Confira se ${API_URL} é mesmo a API do CRM.`);
  }
}
