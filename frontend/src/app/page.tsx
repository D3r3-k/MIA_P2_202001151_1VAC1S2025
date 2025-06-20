"use client";

import CodeEditor from "@/components/CodeEditor/CodeEditor";
import { useMia } from "@/hooks/useMia";
import { TerminalSquare } from "lucide-react";
import { useState } from "react";

export default function HomePage() {
  // Hooks
  const { executeCommand } = useMia();
  // States
  const [consoleInput, setConsoleInput] = useState<string>("");
  const [response, setResponse] = useState<string>("");
  // Effects
  // Handlers
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
    input.onchange = (e: any) => {
      const file = e.target.files[0];
      if (file) {
        const reader = new FileReader();
        reader.onload = (event) => {
          const text = event.target?.result as string;
          setConsoleInput(text);
        };
        reader.readAsText(file);
      }
    };
    input.click();
  }
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
    if (!confirm("¿Estás seguro de que quieres crear un nuevo archivo? Se perderá el contenido actual.")) {
      return;
    }
    setConsoleInput("");
    setResponse("");
  };
  // Renders
  return (
    <main className="flex-1 p-6 ml-72">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-white mb-2">
          Consola de Discos
        </h1>
        <p className="text-gray-400">
          Bienvenido a la consola de discos. Aquí puedes interactuar con el sistema de archivos, ejecutar comandos y gestionar particiones.
        </p>
      </div>
      <div className="mb-6 flex justify-center">
        <nav className="flex space-x-4 bg-gray-800 rounded-lg p-4 shadow-md">
          <button
            onClick={handleFileOpen}
            className="px-4 py-2 bg-gradient-to-r from-gray-500 to-gray-600 text-white rounded-md hover:from-gray-800 hover:to-gray-700 transition-colors font-medium shadow cursor-pointer">
            Abrir Archivo
          </button>
          <button
            onClick={handleFileSave}
            className="px-4 py-2 bg-gradient-to-r from-gray-500 to-gray-600 text-white rounded-md hover:from-gray-800 hover:to-gray-700 transition-colors font-medium shadow cursor-pointer">
            Guardar Archivo
          </button>
          <button
            onClick={handleFileNew}
            className="px-4 py-2 bg-gradient-to-r from-gray-500 to-gray-600 text-white rounded-md hover:from-gray-800 hover:to-gray-700 transition-colors font-medium shadow cursor-pointer">
            Nuevo Archivo
          </button>
        </nav>
      </div>
      <div className="flex space-x-6">
        <div className="flex-1">
          <div className="flex items-center mb-4">
            <TerminalSquare className="text-gray-400 mr-2" />
            <h2 className="text-xl font-semibold text-white">Consola</h2>
          </div>
          <CodeEditor editable onChange={handleConsoleInputChange} value={consoleInput} />
          <div className="mt-4 flex justify-end">
            <button
              onClick={handleExecuteCommand}
              className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 transition-colors cursor-pointer">
              Ejecutar Comando
            </button>
          </div>
        </div>
        <div className="flex-1">
          <div className="flex items-center mb-4">
            <TerminalSquare className="text-gray-400 mr-2" />
            <h2 className="text-xl font-semibold text-white">Salida</h2>
          </div>
          <CodeEditor editable={false} value={response} />
        </div>
      </div>
    </main>
  );
}
