import React from 'react'

export default function SkeletonDriveDisk() {
    return (
        <div
            className="bg-gray-800/50 border-gray-700 hover:bg-gray-800/70 hover:border-corinto-600/70 transition-all duration-200 hover:scale-105 group cursor-pointer border rounded-lg shadow-lg">
            <div className="p-6">
                <div className="flex items-center justify-between mb-4">
                    <div className="relative">
                        <div className="w-14 h-14 bg-gray-600 rounded-full animate-pulse" />
                        <div className="absolute top-0 right-0 w-4 h-4 rounded-full bg-gray-600 animate-pulse" />
                    </div>
                    <div className="text-right">
                        <div className="h-4 bg-gray-600 rounded w-20 mb-2 animate-pulse" />
                    </div>
                </div>

                <div className="space-y-3">
                    <div>
                        <div className="h-5 bg-gray-600 rounded w-40 mb-2 animate-pulse" />
                        <div className="h-3 bg-gray-600 rounded w-48 mb-2 animate-pulse" />
                    </div>

                    <div className="flex justify-between items-center text-sm mt-2">
                        <div className="h-3 bg-gray-600 rounded w-28 animate-pulse" />
                        <div className="h-3 bg-gray-600 rounded w-20 animate-pulse" />
                    </div>

                    <div className="flex items-center justify-between pt-3 border-t border-gray-700 mt-3">
                        <div className="flex items-center space-x-2">
                            <div className="w-2.5 h-2.5 rounded-full bg-gray-600 animate-pulse" />
                            <div className="h-3 bg-gray-600 rounded w-24 animate-pulse" />
                        </div>
                        <div className="h-3 bg-gray-600 rounded w-12 animate-pulse" />
                    </div>
                </div>
            </div>
        </div>
    )
}
