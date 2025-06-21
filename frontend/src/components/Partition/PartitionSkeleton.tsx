export function PartitionSkeleton() {
    return (
        <div className="bg-gray-800/30 border border-gray-700/50 rounded-lg p-6 animate-pulse">
            <div className="grid grid-cols-1 lg:grid-cols-8 gap-6 items-center">
                <div className="lg:col-span-5">
                    <div className="flex items-center space-x-4">
                        <div className="p-3 rounded-xl border bg-gray-700/50 border-gray-600">
                            <div className="w-6 h-6 bg-gray-600 rounded"></div>
                        </div>
                        <div className="flex-1 space-y-2">
                            <div className="h-4 bg-gray-600 rounded w-1/2"></div>
                            <div className="h-3 bg-gray-600 rounded w-1/3"></div>
                            <div className="flex gap-2 mt-1">
                                <div className="h-4 w-16 bg-gray-700 rounded"></div>
                                <div className="h-4 w-12 bg-gray-700 rounded"></div>
                            </div>
                        </div>
                    </div>
                </div>
                <div className="lg:col-span-1 space-y-2">
                    <div className="flex items-center px-3 py-2 rounded-lg border bg-gray-700/30 border-gray-600">
                        <div className="w-4 h-4 mr-2 bg-gray-600 rounded-full"></div>
                        <div className="h-3 bg-gray-600 rounded w-12"></div>
                    </div>
                </div>
                <div className="lg:col-span-2 text-center">
                    <div className="h-3 w-20 bg-gray-600 mx-auto mb-1 rounded"></div>
                    <div className="h-4 w-32 bg-gray-700 mx-auto rounded"></div>
                </div>
            </div>
            <div className="mt-4 pt-4 border-t border-gray-700/30 text-xs text-gray-500 flex justify-between">
                <div className="h-3 w-40 bg-gray-600 rounded"></div>
                <div className="flex gap-4">
                    <div className="h-3 w-24 bg-gray-600 rounded"></div>
                    <div className="h-3 w-24 bg-gray-600 rounded"></div>
                </div>
            </div>
        </div>
    );
}
