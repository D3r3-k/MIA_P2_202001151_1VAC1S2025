import { HardDrive, Activity, Database, HardDriveIcon, TrendingUp, Zap } from 'lucide-react';

interface StatsCardsProps {
    stats: {
        totalDisks: number;
        activeDisks: number;
        totalPartitions: number;
        totalCapacity: string;
    };
}

export function StatsCards({ stats }: StatsCardsProps) {
    const cards = [
        {
            title: 'Total Discos',
            value: stats.totalDisks,
            icon: HardDrive,
            color: 'text-red-400',
            bgColor: 'from-red-500/10 to-red-600/5',
            borderColor: 'border-red-500/20',
            accentColor: 'bg-red-500/20',
            progressColor: 'from-red-500 to-red-400',
            change: '+2',
            changeType: 'positive' as const,
            subtitle: 'Dispositivos registrados',
            progress: 85
        },
        {
            title: 'Discos Activos',
            value: stats.activeDisks,
            icon: Activity,
            color: 'text-green-400',
            bgColor: 'from-green-500/10 to-green-600/5',
            borderColor: 'border-green-500/20',
            accentColor: 'bg-green-500/20',
            progressColor: 'from-green-500 to-green-400',
            change: '+1',
            changeType: 'positive' as const,
            subtitle: 'En funcionamiento',
            progress: 92
        },
        {
            title: 'Particiones',
            value: stats.totalPartitions,
            icon: Database,
            color: 'text-blue-400',
            bgColor: 'from-blue-500/10 to-blue-600/5',
            borderColor: 'border-blue-500/20',
            accentColor: 'bg-blue-500/20',
            progressColor: 'from-blue-500 to-blue-400',
            change: '0',
            changeType: 'neutral' as const,
            subtitle: 'Volúmenes configurados',
            progress: 78
        },
        {
            title: 'Capacidad Total',
            value: stats.totalCapacity,
            icon: HardDriveIcon,
            color: 'text-purple-400',
            bgColor: 'from-purple-500/10 to-purple-600/5',
            borderColor: 'border-purple-500/20',
            accentColor: 'bg-purple-500/20',
            progressColor: 'from-purple-500 to-purple-400',
            change: '500GB',
            changeType: 'positive' as const,
            subtitle: 'Espacio disponible',
            progress: 65
        }
    ];

    const getChangeIcon = (type: 'positive' | 'negative' | 'neutral') => {
        switch (type) {
            case 'positive':
                return <TrendingUp className="w-3 h-3" />;
            case 'negative':
                return <TrendingUp className="w-3 h-3 rotate-180" />;
            default:
                return <Zap className="w-3 h-3" />;
        }
    };

    const getChangeColor = (type: 'positive' | 'negative' | 'neutral') => {
        switch (type) {
            case 'positive':
                return 'text-green-400 bg-green-500/10 border-green-500/20';
            case 'negative':
                return 'text-red-400 bg-red-500/10 border-red-500/20';
            default:
                return 'text-gray-400 bg-gray-500/10 border-gray-500/20';
        }
    };

    return (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {cards.map((card, index) => (
                <Card
                    key={index}
                    className={`bg-gradient-to-br ${card.bgColor} border ${card.borderColor} backdrop-blur-sm relative overflow-hidden bg-gray-800/30`}
                >
                    <CardContent className="p-6 relative z-10">
                        {/* Header with Icon and Change Indicator */}
                        <div className="flex items-start justify-between mb-4">
                            <div className={`p-3 rounded-xl ${card.accentColor} border ${card.borderColor} shadow-lg backdrop-blur-sm`}>
                                <card.icon className={`w-6 h-6 ${card.color}`} />
                            </div>
                            <div className={`flex items-center space-x-1 px-2 py-1 rounded-lg text-xs font-medium border ${getChangeColor(card.changeType)}`}>
                                {getChangeIcon(card.changeType)}
                                <span>{card.change}</span>
                            </div>
                        </div>

                        {/* Main Content */}
                        <div className="space-y-2">
                            <div>
                                <p className="text-gray-400 text-sm font-medium">
                                    {card.title}
                                </p>
                                <p className="text-3xl font-bold text-white tracking-tight">
                                    {card.value}
                                </p>
                            </div>

                            {/* Subtitle */}
                            <p className="text-xs text-gray-500 font-medium">
                                {card.subtitle}
                            </p>
                        </div>

                        {/* Progress Bar */}
                        <div className="mt-4 pt-4 border-t border-gray-700/30">
                            <div className="flex items-center justify-between mb-2">
                                <span className="text-xs text-gray-400">Utilización</span>
                                <span className="text-xs text-gray-300 font-medium">
                                    {card.progress}%
                                </span>
                            </div>
                            <div className="w-full h-1.5 bg-gray-700/50 rounded-full overflow-hidden">
                                <div
                                    className={`h-full bg-gradient-to-r ${card.progressColor} rounded-full transition-all duration-1000 ease-out`}
                                    style={{ width: `${card.progress}%` }}
                                />
                            </div>
                        </div>
                    </CardContent>

                    {/* Subtle Background Pattern for Dark Mode */}
                    <div className="absolute inset-0 opacity-5">
                        <div className="absolute top-0 right-0 w-32 h-32 bg-gradient-to-bl from-gray-300/20 to-transparent rounded-full -translate-y-16 translate-x-16" />
                        <div className="absolute bottom-0 left-0 w-24 h-24 bg-gradient-to-tr from-gray-300/10 to-transparent rounded-full translate-y-12 -translate-x-12" />
                    </div>
                </Card>
            ))}
        </div>
    );
}