interface DriveStatsProps {
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

export function DriveStats({
    title,
    value,
    color: card,
    icon: Icon,
}: DriveStatsProps) {

    return (
        <div className={`bg-gradient-to-br ${card.bgColor} border ${card.borderColor} backdrop-blur-sm relative overflow-hidden bg-gray-800/30 rounded-md`}>
            <div className="p-4 relative z-10 flex flex-col justify-center items-center gap-3">
                <div className="flex gap-4">
                    <div className="flex items-start justify-between">
                        <div className={`p-3 rounded-xl ${card.accentColor} border ${card.borderColor} shadow-lg backdrop-blur-sm`}>
                            <Icon className={`w-6 h-6 ${card.color}`} />
                        </div>
                    </div>
                    <div>
                        <p className="text-gray-400 text-sm font-medium">
                            {title}
                        </p>
                    </div>
                </div>
                <div className="flex">
                    <p className="text-2xl font-bold text-white tracking-tight">
                        {value}
                    </p>
                </div>
            </div>
        </div>
    );
}