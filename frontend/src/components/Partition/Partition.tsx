import { HardDrive, Database, Folder, Activity, CheckCircle, Clock, XCircle, Info, Edit, Mountain } from 'lucide-react';
type PartitionProps = {
    partition: {
        name: string
        disk: string
        size: string
        type: string
        filesystem: string
        mountPoint: string
        usedSpace: number
        freeSpace: string
        status: string
        lastCheck: string
        id: string
    }
}

export function Partition({ partition }: PartitionProps) {
    const getStatusConfig = (status: string) => {
        switch (status) {
            case 'mounted':
                return { color: 'text-green-400', bg: 'bg-green-500/20', border: 'border-green-500/30', icon: CheckCircle, text: 'Montada' };
            case 'unmounted':
                return { color: 'text-yellow-400', bg: 'bg-yellow-500/20', border: 'border-yellow-500/30', icon: Clock, text: 'Desmontada' };
            case 'error':
                return { color: 'text-red-400', bg: 'bg-red-500/20', border: 'border-red-500/30', icon: XCircle, text: 'Error' };
            default:
                return { color: 'text-gray-400', bg: 'bg-gray-500/20', border: 'border-gray-500/30', icon: Info, text: 'Desconocido' };
        }
    };

    const getHealthConfig = (health: string) => {
        switch (health) {
            case 'good':
                return { color: 'text-green-400', bg: 'bg-green-500/10', text: 'Saludable' };
            case 'warning':
                return { color: 'text-yellow-400', bg: 'bg-yellow-500/10', text: 'Advertencia' };
            case 'critical':
                return { color: 'text-red-400', bg: 'bg-red-500/10', text: 'Crítico' };
            default:
                return { color: 'text-gray-400', bg: 'bg-gray-500/10', text: 'Desconocido' };
        }
    };

    const getTypeConfig = (type: string) => {
        switch (type.toLowerCase()) {
            case 'primaria':
                return { icon: HardDrive, color: 'text-blue-400', bg: 'bg-blue-500/20', border: 'border-blue-500/30' };
            case 'extendida':
                return { icon: Database, color: 'text-purple-400', bg: 'bg-purple-500/20', border: 'border-purple-500/30' };
            case 'lógica':
                return { icon: Folder, color: 'text-green-400', bg: 'bg-green-500/20', border: 'border-green-500/30' };
            case 'swap':
                return { icon: Activity, color: 'text-orange-400', bg: 'bg-orange-500/20', border: 'border-orange-500/30' };
            default:
                return { icon: Database, color: 'text-gray-400', bg: 'bg-gray-500/20', border: 'border-gray-500/30' };
        }
    };

    const statusConfig = getStatusConfig(partition.status);
    const healthConfig = getHealthConfig("good");
    const typeConfig = getTypeConfig(partition.type);
    const StatusIcon = statusConfig.icon;
    const TypeIcon = typeConfig.icon;

    return (
        <div className="bg-gray-800/30 border border-gray-700/50 rounded-lg p-6 hover:bg-gray-800/50 transition-all duration-200">
            <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-center">
                {/* Información principal */}
                <div className="lg:col-span-5">
                    <div className="flex items-center space-x-4">
                        <div className={`p-3 rounded-xl border ${typeConfig.bg} ${typeConfig.border}`}>
                            <TypeIcon className={`w-6 h-6 ${typeConfig.color}`} />
                        </div>
                        <div className="flex-1">
                            <h3 className="text-lg font-semibold text-white">{partition.name}</h3>
                            <p className="text-sm text-gray-400">{partition.disk} • {partition.size}</p>
                            <div className="flex gap-2 mt-1">
                                <span className={`text-xs px-2 py-1 rounded ${typeConfig.bg} ${typeConfig.border} ${typeConfig.color}`}>{partition.type}</span>
                                <span className="text-xs px-2 py-1 rounded text-gray-400 bg-gray-700/30 border border-gray-600">{partition.filesystem}</span>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Estado y Salud */}
                <div className="lg:col-span-2 space-y-2">
                    <div className={`flex items-center px-3 py-2 rounded-lg border ${statusConfig.bg} ${statusConfig.border}`}>
                        <StatusIcon className={`w-4 h-4 mr-2 ${statusConfig.color}`} />
                        <span className={`text-sm font-medium ${statusConfig.color}`}>{statusConfig.text}</span>
                    </div>
                </div>

                {/* Punto de montaje */}
                <div className="lg:col-span-2 text-center">
                    <p className="text-xs text-gray-400 mb-1">Punto de montaje</p>
                    <code className="text-sm text-white bg-gray-900/50 px-2 py-1 rounded">{partition.mountPoint}</code>
                </div>

                {/* Uso de espacio */}
                <div className="lg:col-span-3 space-y-2">
                    <div className="flex justify-between text-sm text-gray-400">
                        <span>Uso del espacio</span>
                        <span className="text-white font-medium">{partition.usedSpace}%</span>
                    </div>
                    <div className="h-2 w-full bg-gray-700/50 rounded">
                        <div className="h-full bg-green-500 rounded" style={{ width: `${partition.usedSpace}%` }} />
                    </div>
                    <div className="flex justify-between text-xs text-gray-500">
                        <span>Libre: {partition.freeSpace}</span>
                        <span>Total: {partition.size}</span>
                    </div>
                </div>
            </div>

            {/* Información adicional */}
            <div className="mt-4 pt-4 border-t border-gray-700/30 text-xs text-gray-500 flex justify-between">
                <span>Última verificación: {partition.lastCheck}</span>
                <div className="flex gap-4">
                    <span>ID: {partition.id}</span>
                    <span>Sector: {Math.floor(Math.random() * 1000000)}</span>
                </div>
            </div>
        </div>
    );
}
