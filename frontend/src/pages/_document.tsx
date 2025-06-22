import { Html, Head, Main, NextScript } from "next/document";

export default function Document() {
  return (
    <Html lang="en">
      <Head />
      <body className="antialiased">
        <div className="min-h-screen bg-gray-950">
          <Main />
          <NextScript />
        </div>
      </body>
    </Html>
  );
}
