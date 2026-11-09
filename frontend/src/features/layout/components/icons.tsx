import type { ComponentProps } from "react";

function Icone(props: ComponentProps<"svg">) {
  return (
    <svg
      viewBox="0 0 24 24"
      width="22"
      height="22"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.75"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
      {...props}
    />
  );
}

export const IconeInicio = (p: ComponentProps<"svg">) => (
  <Icone {...p}>
    <path d="M3 10.5 12 3l9 7.5" />
    <path d="M5 9.5V20h14V9.5" />
    <path d="M10 20v-5.5h4V20" />
  </Icone>
);

export const IconeCrm = (p: ComponentProps<"svg">) => (
  <Icone {...p}>
    <circle cx="9" cy="8" r="3.25" />
    <path d="M3 19.5c.5-3.2 3-5 6-5s5.5 1.8 6 5" />
    <path d="M16 5.2a3.25 3.25 0 0 1 0 5.6" />
    <path d="M18 14.8c1.7.6 2.7 2.1 3 4.7" />
  </Icone>
);

export const IconeConfiguracoes = (p: ComponentProps<"svg">) => (
  <Icone {...p}>
    <circle cx="12" cy="12" r="3" />
    <path d="M12 2.5v2.4M12 19.1v2.4M4.6 4.6l1.7 1.7M17.7 17.7l1.7 1.7M2.5 12h2.4M19.1 12h2.4M4.6 19.4l1.7-1.7M17.7 6.3l1.7-1.7" />
  </Icone>
);

export const IconeMenu = (p: ComponentProps<"svg">) => (
  <Icone {...p}>
    <path d="M4 7h16M4 12h16M4 17h16" />
  </Icone>
);

export const IconeFechar = (p: ComponentProps<"svg">) => (
  <Icone {...p}>
    <path d="M6 6l12 12M18 6 6 18" />
  </Icone>
);

export const IconeSetaEsquerda = (p: ComponentProps<"svg">) => (
  <Icone {...p}>
    <path d="m14 7-5 5 5 5" />
  </Icone>
);

export const IconeSetaDireita = (p: ComponentProps<"svg">) => (
  <Icone {...p}>
    <path d="m10 7 5 5-5 5" />
  </Icone>
);
