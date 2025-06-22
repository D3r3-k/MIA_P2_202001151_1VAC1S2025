"use client";

import useFetchs from "@/hooks/useFetchs";
import { FileSystemItemType } from "@/types/GlobalTypes";
import { ArrowLeft, ChevronRight, File, Folder, HardDrive, Maximize2, Minimize2, X } from "lucide-react";
import { useEffect, useState } from "react";


export default function ContentPartition() {
    // Hooks
    const { getFileSystemItems, getContentFile } = useFetchs();
    // States
    const [rootItems, setRootItems] = useState<FileSystemItemType[]>([]);
    const [currentPath, setCurrentPath] = useState('/');
    const [pathHistory, setPathHistory] = useState<string[]>(["/"]);
    const [viewingFile, setViewingFile] = useState<FileSystemItemType | null>(null);
    const [isFileViewMaximized, setIsFileViewMaximized] = useState(false);
    // Effects
    useEffect(() => {
        const fetchFileSystemItems = async () => {
            try {
                const items = await getFileSystemItems('/');
                let root: FileSystemItemType[] = [];
                if (Array.isArray(items)) {
                    root = items;
                } else if ('Root' in items && Array.isArray(items.Root)) {
                    root = items.Root;
                }
                const sorted = root.sort((a: FileSystemItemType, b: FileSystemItemType) => {
                    if (a.Type === 'folder' && b.Type !== 'folder') return -1;
                    if (a.Type !== 'folder' && b.Type === 'folder') return 1;
                    return a.Name.localeCompare(b.Name);
                });
                setRootItems(sorted);
            } catch (error) {
                console.error("Error fetching file system items:", error);
                setRootItems([]);
            }
        };

        fetchFileSystemItems();
    }, []);

    // Handlers
    const viewFile = async (fileName: string) => {
        // Busca un archivo por nombre en el path actual
        const findFileByName = (items: FileSystemItemType[], path: string, name: string): FileSystemItemType | null => {
            const segments = path.split('/').filter(Boolean);
            let currentItems = items;

            for (const segment of segments) {
                const folder = currentItems.find(item => item.Name === segment && item.Type === 'folder');
                if (folder && folder.Children) {
                    currentItems = folder.Children;
                } else {
                    return null;
                }
            }

            return currentItems.find(item => item.Name === name && item.Type === 'file') || null;
        };

        const file = findFileByName(rootItems, currentPath, fileName);
        if (!file) {
            console.warn(`Archivo "${fileName}" no encontrado en ${currentPath}`);
            return;
        }

        try {
            const data = await getContentFile(file.Path);
            if (!data || typeof data !== 'object') {
                console.error("Datos del archivo no válidos:", data);
                return;
            }
            const dataP = data as FileSystemItemType;
            const fileWithContent: FileSystemItemType = {
                ...file,
                Content: dataP.Content || "",
                CreatedAt: dataP.CreatedAt || "",
                Extension: dataP.Extension || "",
                Owner: dataP.Owner || "",
                Permissions: dataP.Permissions || "",
                Size: dataP.Size || "",
            };
            setViewingFile(fileWithContent);
            setIsFileViewMaximized(false);
        } catch (error) {
            console.error("Error al obtener contenido del archivo:", error);
        }
    };



    const closeFileView = () => {
        setViewingFile(null);
        setIsFileViewMaximized(false);
    };
    // Functions
    const getFileIcon = (item: FileSystemItemType) => {
        if (item.Type === 'folder') return Folder;
        return File;
    }
    const navigateToFolder = (folderName: string) => {
        const newPath = currentPath === '/' ? `/${folderName}` : `${currentPath}/${folderName}`;
        setCurrentPath(newPath);
        setPathHistory([...pathHistory, newPath]);
    };
    const goBack = () => {
        if (pathHistory.length > 1) {
            const newPathHistory = [...pathHistory];
            newPathHistory.pop();
            setPathHistory(newPathHistory);
            setCurrentPath(newPathHistory[newPathHistory.length - 1] || '/');
        }
    };
    const getPathSegments = () => {
        if (currentPath === '/') return [{ name: '/', path: '/' }];
        const segments = currentPath.split('/').filter(Boolean);
        const result = [{ name: '/', path: '/' }];
        let currentSegmentPath = '';
        segments.forEach(segment => {
            currentSegmentPath += '/' + segment;
            result.push({ name: segment, path: currentSegmentPath });
        });
        return result;
    };
    // Renders
    const getItemsByPath = (path: string, items: FileSystemItemType[]): FileSystemItemType[] => {
        if (path === '/' || path === '') {
            // Ordenar primero carpetas, luego archivos
            return [...items].sort((a, b) => {
                if (a.Type === 'folder' && b.Type !== 'folder') return -1;
                if (a.Type !== 'folder' && b.Type === 'folder') return 1;
                return a.Name.localeCompare(b.Name);
            });
        }
        const segments = path.split('/').filter(Boolean);
        let currentItems = items;
        for (const segment of segments) {
            const found = currentItems.find(item => item.Name === segment && item.Type === 'folder');
            if (found && found.Children) {
                currentItems = found.Children;
            } else {
                currentItems = [];
                break;
            }
        }
        return [...currentItems].sort((a, b) => {
            if (a.Type === 'folder' && b.Type !== 'folder') return -1;
            if (a.Type !== 'folder' && b.Type === 'folder') return 1;
            return a.Name.localeCompare(b.Name);
        });
    };

    const filteredItems = getItemsByPath(currentPath, rootItems);
    return (
        <>
            {/* Barra de navegación */}
            <div className="bg-gray-800/30 border border-gray-700/50 rounded-lg">
                <div className="p-4">
                    <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center space-x-2">
                            <button
                                onClick={goBack}
                                disabled={pathHistory.length <= 1}
                                className="bg-gray-700/50 border border-gray-600 text-gray-300 hover:bg-gray-600/50 disabled:opacity-50 rounded text-sm flex items-center px-3 py-2 cursor-pointer"
                            >
                                <ArrowLeft className="w-4 h-4" />
                            </button>
                            <div className="flex items-center space-x-1 bg-gray-900/50 px-3 py-2 rounded-lg">
                                {getPathSegments().map((segment, index) => (
                                    <div key={index} className="flex items-center">
                                        {index > 0 && (
                                            <ChevronRight className="w-4 h-4 text-gray-500 mx-1" />
                                        )}
                                        {index === getPathSegments().length - 1 ? (
                                            <span className="text-sm text-white font-semibold px-2 py-1 rounded bg-gray-700/70">
                                                {segment.name}
                                            </span>
                                        ) : (
                                            <button
                                                onClick={() => {
                                                    const newPath = segment.path;
                                                    setCurrentPath(newPath);
                                                    setPathHistory(pathHistory.slice(0, pathHistory.indexOf(newPath) + 1));
                                                }}
                                                className="text-sm text-gray-300 hover:text-white transition-colors px-2 py-1 rounded hover:bg-gray-700/50 cursor-pointer"
                                            >
                                                {segment.name}
                                            </button>
                                        )}
                                    </div>
                                ))}
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Contenido del directorio */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <div className={`bg-gray-800/30 border border-gray-700/50 rounded-lg transition-all duration-300 ${viewingFile ? "lg:col-span-2" : "lg:col-span-3"}`}>
                    {/* Header */}
                    <div className="p-4 border-b border-gray-700/50">
                        <div className="text-white flex items-center gap-2 text-lg font-semibold">
                            <HardDrive className="w-5 h-5 text-blue-400" />
                            Contenido de {currentPath === '/' ? 'Raíz' : currentPath.split('/').pop()}
                        </div>
                    </div>
                    {/* Content */}
                    <div className="p-4">
                        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
                            {filteredItems.map((item) => {
                                const ItemIcon = getFileIcon(item);
                                return (
                                    <div
                                        key={item.ID}
                                        className="p-4 rounded-lg border transition-all duration-200 hover:scale-105 bg-gray-700/30 border-gray-600/50 hover:bg-gray-700/50 cursor-pointer"
                                        onClick={() => {
                                            if (item.Type === "folder") {
                                                navigateToFolder(item.Name);
                                            } else {
                                                viewFile(item.Name);
                                            }
                                        }}
                                    >
                                        <div className="text-center space-y-2">
                                            {ItemIcon ? (
                                                <ItemIcon className={`w-8 h-8 mx-auto ${item.Type === "folder" ? "text-blue-400" : "text-corinto-400"}`} />
                                            ) : (
                                                <span className="w-8 h-8 mx-auto inline-block" />
                                            )}
                                            <div>
                                                <p
                                                    className="text-white text-sm font-medium truncate"
                                                    title={item.Name}
                                                >
                                                    {item.Name}
                                                </p>
                                                {item.Size && (
                                                    <p className="text-gray-400 text-xs">{item.Size}</p>
                                                )}
                                            </div>
                                        </div>
                                    </div>
                                );
                            })}
                        </div>
                    </div>
                </div>

                {/* Vista de archivo */}
                {viewingFile && (
                    <div className={`transition-all duration-300 ${isFileViewMaximized ? "fixed inset-0 z-50 bg-gray-950 p-6" : "lg:col-span-1"}`}>
                        <div className="bg-gray-800/30 border border-gray-700/50 h-full flex flex-col rounded-lg overflow-hidden">
                            {/* Header */}
                            <div className="pb-3 p-4 border-b border-gray-700/30 flex-shrink-0">
                                <div className="flex items-center justify-between">
                                    <div className="flex items-center gap-2 truncate text-white font-semibold text-base">
                                        {(() => {
                                            const FileIcon = getFileIcon(viewingFile);
                                            return (
                                                <FileIcon className={`w-5 h-5`} />
                                            );
                                        })()}
                                        <span className="truncate">{viewingFile.Name}</span>
                                    </div>
                                    <div className="flex items-center space-x-1 flex-shrink-0">
                                        <button
                                            onClick={() => setIsFileViewMaximized(!isFileViewMaximized)}
                                            className="text-gray-400 hover:text-white p-2 rounded cursor-pointer"
                                        >
                                            {isFileViewMaximized ? (
                                                <Minimize2 className="w-4 h-4" />
                                            ) : (
                                                <Maximize2 className="w-4 h-4" />
                                            )}
                                        </button>
                                        <button
                                            onClick={closeFileView}
                                            className="text-gray-400 hover:text-white p-2 rounded cursor-pointer"
                                        >
                                            <X className="w-4 h-4" />
                                        </button>
                                    </div>
                                </div>
                                <div className="flex items-center space-x-4 text-xs text-gray-400 mt-2">
                                    <span>{viewingFile.Size}</span>
                                    <span>{viewingFile.CreatedAt}</span>
                                    <span>{viewingFile.Owner}</span>
                                </div>
                            </div>
                            {/* Content */}
                            <div className="flex-1 flex flex-col min-h-0 p-4">
                                <div className="flex-1 min-h-0">
                                    <textarea
                                        readOnly
                                        value={viewingFile.Content || "No hay contenido disponible para mostrar."}
                                        className="w-full h-full min-h-[400px] bg-black/50 border border-gray-600 text-green-400 font-mono text-sm resize-none rounded p-3 focus:ring-0 focus:outline-none"
                                    />
                                </div>
                                {/* Footer */}
                                <div className="mt-4 pt-4 border-t border-gray-700/30 flex-shrink-0">
                                    <div className="flex items-center justify-between">
                                        <div className="flex items-center space-x-2">
                                            <span className="px-2 py-1 rounded border border-gray-600 bg-gray-700/30 text-xs text-gray-400">
                                                {viewingFile.Extension?.toUpperCase() || "ARCHIVO"}
                                            </span>
                                            <span className="px-2 py-1 rounded border border-gray-600 bg-gray-700/30 text-xs text-gray-400">
                                                {viewingFile.Permissions}
                                            </span>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </>

    )
}
