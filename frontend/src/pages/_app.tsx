"use client";

import Loading from "@/components/Loading/Loading";
import Sidebar from "@/components/Sidebar/Sidebar";
import Toast from "@/components/Toast/Toast";
import { MiaProvider } from "@/contexts/MiaContext";
import "@/styles/globals.css";
import type { AppProps } from "next/app";

export default function App({ Component, pageProps }: AppProps) {
  return (
    <MiaProvider>
      <Sidebar />
      <div className="relative">
        <Component {...pageProps} />
        <Toast />
        <Loading />
      </div>
    </MiaProvider>
  );
}
