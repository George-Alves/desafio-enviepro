import type { Metadata } from "next";
import { Bricolage_Grotesque, Figtree } from "next/font/google";
import { AppShell } from "@/features/layout";
import "./globals.css";

const titulo = Bricolage_Grotesque({
  variable: "--font-titulo",
  subsets: ["latin"],
});

const corpo = Figtree({
  variable: "--font-corpo",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "EnviePro CRM",
  description: "Clientes, funil de vendas e resumo de vendas em um só painel.",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="pt-BR" className={`${titulo.variable} ${corpo.variable} antialiased`}>
      <body>
        <AppShell>{children}</AppShell>
      </body>
    </html>
  );
}
