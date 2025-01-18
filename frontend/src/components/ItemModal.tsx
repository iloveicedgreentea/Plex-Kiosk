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
        <div className="aspect-[2/3] relative">
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

    console.log('Rendering modal with item:', item); // Debug log

    return (
        <Dialog open={isOpen} onOpenChange={() => onClose()}>
            <DialogContent className="bg-[#2d2d2d] text-white p-6 max-w-4xl w-full">
                <div className="flex justify-between items-start">
                    <h2 className="text-2xl font-bold text-[#e5a00d]">{item.title}</h2>
                    <button
                        onClick={() => onClose()}
                        className="text-gray-400 hover:text-white"
                    >
                        <X size={24} />
                    </button>
                </div>

                <div className="mt-6 grid grid-cols-1 md:grid-cols-2 gap-6">
                    <DebugImage src={item.thumb_url || ''} alt={item.title} />

                    <div className="space-y-6">
                        <div>
                            <h3 className="text-[#e5a00d] font-semibold mb-2">Details</h3>
                            <div className="space-y-1">
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
                                <p className="text-gray-300 leading-relaxed">{item.description}</p>
                            </div>
                        )}

                        {item.trailer_url && (
                            <div>
                                <h3 className="text-[#e5a00d] font-semibold mb-2">Trailer</h3>
                                <a
                                    href={item.trailer_url}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="text-blue-400 hover:text-blue-300"
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