"use client";

import { useState } from "react";
import type { AppProps } from "next/app";
import "@/styles/globals.css";
import { MiaProvider } from "@/contexts/MiaContext";
import Sidebar from "@/components/Sidebar/Sidebar";
import Toast from "@/components/Toast/Toast";
import Loading from "@/components/Loading/Loading";
import { Route } from "@/types/GlobalTypes";

export default function App({ Component, pageProps }: AppProps) {
  const [route, setRoute] = useState<Route>("/");
  
  return (
    <MiaProvider>
      <Sidebar activeRoute={route} setRoute={setRoute} />
      <div className="relative">
        <Component {...pageProps} route={route} setRoute={setRoute} />
        <Toast />
        <Loading />
      </div>
    </MiaProvider>
  );
}
