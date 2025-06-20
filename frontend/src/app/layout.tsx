import type { Metadata } from "next";
import "@/styles/globals.css";
import { Inter } from 'next/font/google';
import { MiaProvider } from "@/contexts/MiaContext";
import Loading from "@/components/Loading/Loading";
import Sidebar from "@/components/Sidebar/Sidebar";


const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: "Sistema de Archivos - MIA F2",
  description: "Segunda fase del proyecto 1 de Manejo e implementación de Archivos",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="es">
      <body
        className={`${inter.className}`}
      >
        <MiaProvider>
          <div className="min-h-screen bg-gray-950">
            <div className="flex">
              <Sidebar />
              {children}
              <Loading />
            </div>
          </div>
        </MiaProvider>
      </body>
    </html>
  );
}
