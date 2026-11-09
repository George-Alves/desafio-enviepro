import type { ComponentProps } from "react";

const TONS = {
  neutro: "",
  violeta: "vidro-violeta",
  lavanda: "vidro-lavanda",
  indigo: "vidro-indigo",
  orquidea: "vidro-orquidea",
} as const;

const FORMATOS = {
  cartao: "rounded-[1.75rem]",
  pilula: "rounded-full",
} as const;

type Props = ComponentProps<"div"> & {
  /** Cor do vidro. */
  tom?: keyof typeof TONS;
  formato?: keyof typeof FORMATOS;
  /** O cartão sobe um pouco ao passar o mouse. */
  eleva?: boolean;
};

/** Painel de vidro (glassmorphism). Toda tela usa este componente em vez de repetir o estilo. */
export function GlassCard({
  tom = "neutro",
  formato = "cartao",
  eleva = false,
  className = "",
  ...props
}: Props) {
  const classes = ["vidro", TONS[tom], FORMATOS[formato], eleva ? "vidro-eleva" : "", className]
    .filter(Boolean)
    .join(" ");
  return <div className={classes} {...props} />;
}
