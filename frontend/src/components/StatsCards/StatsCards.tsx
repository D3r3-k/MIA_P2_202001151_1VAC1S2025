import { HardDrive, Activity, Database, HardDriveIcon, TrendingUp, Zap } from 'lucide-react';

interface StatsCardsProps {
    title: string;
    value: number | string;
    color: {
        color: string;
        bgColor: string;
        borderColor: string;
        accentColor: string;
    }
    icon: React.ComponentType<any>;
}

export function StatsCards({
    title,
    value,
    color: card,
    icon: Icon,
}: StatsCardsProps) {

    return (
        <div className={`bg-gradient-to-br ${card.bgColor} border ${card.borderColor} backdrop-blur-sm relative overflow-hidden bg-gray-800/30 rounded-md`}>
                <div className="p-6 relative z-10 flex gap-6">
                    <div className="flex items-start justify-between mb-4">
                        <div className={`p-3 rounded-xl ${card.accentColor} border ${card.borderColor} shadow-lg backdrop-blur-sm`}>
                            <Icon className={`w-6 h-6 ${card.color}`} />
                        </div>
                    </div>
                    <div className="space-y-2">
                        <div>
                            <p className="text-gray-400 text-sm font-medium">
                                {title}
                            </p>
                            <p className="text-3xl font-bold text-white tracking-tight">
                                {value}
                            </p>
                        </div>
                    </div>
                </div>
            </div>
    );
}