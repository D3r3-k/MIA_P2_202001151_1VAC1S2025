import { StreamLanguage } from "@codemirror/language";
import ReactCodeMirror from "@uiw/react-codemirror";
import { python } from "@codemirror/legacy-modes/mode/python";
import { vscodeDark } from '@uiw/codemirror-theme-vscode';

interface CodeEditorProps {
    editable?: boolean;
    value?: string;
    onChange?: (value: string) => void;
    size?: "small" | "medium" | "large";
}

export default function CodeEditor({ editable, value, onChange, size }: CodeEditorProps) {
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
            height={size === "small" ? "200px" : size === "medium" ? "300px" : "400px"}
            theme={vscodeDark}
            editable={editable}
            readOnly={!editable}
            extensions={[StreamLanguage.define(python)]}
            className="text-sm scrollbar-thin dark:scrollbar-thin flex-1 rounded-lg"
            placeholder={editable ? "Escribe tu código aquí..." : "Salida del programa..."}
        />
    )
}
