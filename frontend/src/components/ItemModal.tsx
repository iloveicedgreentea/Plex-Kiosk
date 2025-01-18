import { Dialog, DialogContent } from "@/components/ui/dialog"
import { X } from "lucide-react"
import { LibraryItem } from "@/types"
import { useState } from "react"

interface ItemModalProps {
    item: LibraryItem | null;
    isOpen: boolean;
    onClose: () => void;
}

const formatDate = (dateString: string) => {
    try {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    } catch (error) {
        console.error('Error formatting date:', error, dateString);
        return 'Unknown date';
    }
};

const DebugImage = ({ src, alt }: { src: string; alt: string }) => {
    const [error, setError] = useState(false);

    return (
        <div className="w-full aspect-[2/3] relative">
            <img
                src={src}
                alt={alt}
                className={`w-full h-full object-cover rounded-lg ${error ? 'hidden' : ''}`}
                onError={() => {
                    console.error('Image failed to load:', src);
                    setError(true);
                }}
                onLoad={() => console.log('Image loaded successfully:', src)}
            />
            {error && (
                <div className="w-full h-full bg-[#1a1a1a] flex items-center justify-center text-gray-500 rounded-lg">
                    Failed to load image
                </div>
            )}
        </div>
    );
};

export function ItemModal({ item, isOpen, onClose }: ItemModalProps) {
    if (!item) return null;

    return (
        <Dialog open={isOpen} onOpenChange={() => onClose()}>
            <DialogContent className="bg-[#2d2d2d] text-white w-[calc(100%-2rem)] max-w-3xl mx-auto p-4 overflow-y-auto max-h-[90vh] rounded-lg">
                <div className="flex justify-between items-start mb-4">
                    <h2 className="text-xl md:text-2xl font-bold text-[#e5a00d] pr-8">{item.title}</h2>
                    <button
                        onClick={() => onClose()}
                        className="text-gray-400 hover:text-white shrink-0"
                    >
                        <X size={24} />
                    </button>
                </div>

                <div className="flex flex-col md:grid md:grid-cols-2 gap-6">
                    <div className="w-full max-w-sm mx-auto md:max-w-none">
                        <DebugImage src={item.thumb_url || ''} alt={item.title} />
                    </div>

                    <div className="space-y-4">
                        <div>
                            <h3 className="text-[#e5a00d] font-semibold mb-2">Details</h3>
                            <div className="space-y-1 text-sm md:text-base">
                                <p className="text-gray-300">Year: {item.year || 'N/A'}</p>
                                <p className="text-gray-300">
                                    Rating: {item.rating ? `${(item.rating * 10).toFixed(1)}%` : 'N/A'}
                                </p>
                                <p className="text-gray-300">
                                    Added: {formatDate(item.added_at)}
                                </p>
                            </div>
                        </div>

                        {item.cast && item.cast.length > 0 && (
                            <div>
                                <h3 className="text-[#e5a00d] font-semibold mb-2">Cast</h3>
                                <div className="flex flex-wrap gap-2">
                                    {item.cast.slice(0, 5).map((actor, index) => (
                                        <span
                                            key={index}
                                            className="bg-[#1a1a1a] px-2 py-1 rounded text-sm"
                                        >
                                            {actor}
                                        </span>
                                    ))}
                                </div>
                            </div>
                        )}

                        {item.description && (
                            <div>
                                <h3 className="text-[#e5a00d] font-semibold mb-2">Description</h3>
                                <p className="text-gray-300 text-sm md:text-base leading-relaxed">{item.description}</p>
                            </div>
                        )}

                        {item.trailer_url && (
                            <div>
                                <h3 className="text-[#e5a00d] font-semibold mb-2">Trailer</h3>
                                <a
                                    href={item.trailer_url}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="text-blue-400 hover:text-blue-300 text-sm md:text-base"
                                >
                                    Watch Trailer
                                </a>
                            </div>
                        )}
                    </div>
                </div>
            </DialogContent>
        </Dialog>
    );
}