export default function SkeletonStatsCard() {
    return (
        <div className="bg-gradient-to-br from-gray-700 to-gray-800 border border-gray-700 backdrop-blur-sm relative overflow-hidden bg-gray-900/30 rounded-md animate-pulse pb-3">
            <div className="p-6 relative z-10 flex gap-6">
                <div className="flex items-start justify-between mb-4">
                    <div className="p-3 rounded-xl bg-gray-800 border border-gray-700 shadow-lg backdrop-blur-sm">
                        <div className="w-6 h-6 bg-gray-600 rounded"></div>
                    </div>
                </div>
                <div className="space-y-2">
                    <div>
                        <div className="h-4 w-24 bg-gray-700 rounded mb-3"></div>
                        <div className="h-8 w-16 bg-gray-700 rounded"></div>
                    </div>
                </div>
            </div>
        </div>
    )
}
