import { StreamLanguage } from "@codemirror/language";
import ReactCodeMirror from "@uiw/react-codemirror";
import { python } from "@codemirror/legacy-modes/mode/python";
import { vscodeDark } from '@uiw/codemirror-theme-vscode';

interface CodeEditorProps {
    editable?: boolean;
    value?: string;
    onChange?: (value: string) => void;
}

export default function CodeEditor({ editable, value, onChange }: CodeEditorProps) {
    // Hooks
    // States
    // Effects
    // Handlers
    // Functions
    // Renders
    return (
        <ReactCodeMirror
            value={value}
            onChange={onChange}
            height="600px"
            theme={vscodeDark}
            editable={editable}
            extensions={[StreamLanguage.define(python)]}
            className="text-sm scrollbar-thin dark:scrollbar-thin"
            placeholder={editable ? "Escribe tu código aquí..." : "Salida del programa..."}
        />
    )
}
