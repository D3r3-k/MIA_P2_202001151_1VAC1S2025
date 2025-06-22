import Sidebar from "@/components/Sidebar/Sidebar";
import { MiaProvider } from "@/contexts/MiaContext";
import "@/styles/globals.css";
import type { AppProps } from "next/app";

export default function App({ Component, pageProps }: AppProps) {
  return (
    <MiaProvider>
      <Sidebar />
      <Component {...pageProps} />
    </MiaProvider>
  );
}
