import { useState } from "react";
import { FileText, FolderOpen, Save, Terminal, Trash2, Play } from "lucide-react";
import Head from "next/head";
import useFetchs from "@/hooks/useFetchs";
import CodeEditor from "@/components/CodeEditor/CodeEditor";

export default function ConsoleView() {
  const { executeCommand } = useFetchs();

  const [consoleInput, setConsoleInput] = useState<string>("");
  const [response, setResponse] = useState<string>("");

  const handleConsoleInputChange = (value: string) => {
    setConsoleInput(value);
  };

  const handleExecuteCommand = async () => {
    const output = await executeCommand(consoleInput);
    setResponse(output);
  };

  const handleFileOpen = () => {
    const input = document.createElement("input");
    input.type = "file";
    input.accept = ".sdaa";
    input.onchange = (e: Event) => {
      const files = (e.target as HTMLInputElement).files;
      if (files && files[0]) {
        const file = files[0];
        const reader = new FileReader();
        reader.onload = (event) => {
          const text = event.target?.result as string;
          setConsoleInput(text);
        };
        reader.readAsText(file);
      }
    };
    input.click();
  };

  const handleFileSave = () => {
    if (!consoleInput) {
      alert("No hay contenido para guardar.");
      return;
    }
    const blob = new Blob([consoleInput], { type: "text/plain" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "archivo.sdaa";
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
  };

  const handleFileNew = () => {
    if (!consoleInput && !response) {
      setConsoleInput("");
      setResponse("");
      return;
    }
    if (!confirm("¿Crear nuevo archivo? Se perderá el contenido actual.")) {
      return;
    }
    setConsoleInput("");
    setResponse("");
  };

  const handleClearConsole = () => {
    setConsoleInput("");
    setResponse("");
  };

  return (
    <>
      <Head>
        <title>Consola de Discos - F2 MIA</title>
        <meta name="description" content="Consola de discos para interactuar con el sistema de archivos y ejecutar comandos." />
        <link rel="icon" href="/favicon.ico" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
      </Head>
      <main className="flex-1 p-6 ml-72 relative">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-white mb-2">
            Consola de Discos
          </h1>
          <p className="text-gray-400">
            Bienvenido a la consola de discos. Aquí puedes interactuar con el sistema de archivos, ejecutar comandos y gestionar particiones.
          </p>
        </div>
        {/* Toolbar */}
        <div className="bg-gray-800/30 border border-gray-700/50 rounded mb-6">
          <div className="p-4">
            <div className="flex flex-wrap gap-3">
              <button
                onClick={handleFileOpen}
                className="flex items-center px-4 py-2 rounded bg-gray-700/50 border border-gray-600 text-gray-300 hover:bg-gray-600/50 hover:text-white transition-colors cursor-pointer"
              >
                <FolderOpen className="w-4 h-4 mr-2" />
                Abrir Archivo
              </button>

              <button
                onClick={handleFileSave}
                className="flex items-center px-4 py-2 rounded bg-gray-700/50 border border-gray-600 text-gray-300 hover:bg-gray-600/50 hover:text-white transition-colors cursor-pointer"
              >
                <Save className="w-4 h-4 mr-2" />
                Guardar Salida
              </button>

              <button
                onClick={handleFileNew}
                className="flex items-center px-4 py-2 rounded bg-gray-700/50 border border-gray-600 text-gray-300 hover:bg-gray-600/50 hover:text-white transition-colors cursor-pointer"
              >
                <FileText className="w-4 h-4 mr-2" />
                Nueva Sesión
              </button>

              <button
                onClick={handleClearConsole}
                className="flex items-center px-4 py-2 rounded bg-red-900/20 border border-red-700/50 text-red-400 hover:bg-red-800/30 hover:text-red-300 transition-colors cursor-pointer"
              >
                <Trash2 className="w-4 h-4 mr-2" />
                Limpiar
              </button>
            </div>
          </div>
        </div>
        <div className="flex flex-col gap-4">
          <div className="bg-gray-900/50 border-gray-700/50 border rounded-lg p-4 flex-1">
            <div className="pb-3">
              <h2 className="text-white flex items-center gap-2">
                <Terminal className="w-5 h-5 text-corinto-400" />
                Salida de Consola
              </h2>
            </div>
            <div className="p-0">
              <CodeEditor editable={false} value={response} size="large" />
            </div>
          </div>
          <div className="bg-gray-900/50 border-gray-700/50 border rounded-lg p-4 flex-1">
            <div className="pb-3">
              <h2 className="text-white flex items-center gap-2">
                <Terminal className="w-5 h-5 text-blue-400" />
                Entrada de Comandos
              </h2>
            </div>
            <div className="p-0 flex gap-4">
              <CodeEditor editable onChange={handleConsoleInputChange} value={consoleInput} size="small" />
              <div className="mt-2">
                <button
                  onClick={handleExecuteCommand}
                  className="flex items-center gap-2 px-5 py-2 rounded-lg bg-gradient-to-r from-green-500 to-green-700 hover:from-green-600 hover:to-green-800 text-white font-semibold shadow-md transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-green-400 disabled:opacity-50 disabled:hover:bg-gradient-to-r disabled:from-green-500 disabled:to-green-700"
                  title="Ejecutar Comando"
                  aria-label="Ejecutar Comando"
                  disabled={!consoleInput.trim()}
                >
                  <Play className="w-5 h-5" />
                  Ejecutar
                </button>
              </div>
            </div>
          </div>
        </div>
      </main>
    </>
  );
}
