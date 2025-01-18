import { useState } from "react"
import { LibraryItem } from "@/types"
import { ItemModal } from "./ItemModal"

interface LibrarySectionProps {
    name: string;
    items: LibraryItem[];
}

export function LibrarySection({ name, items }: LibrarySectionProps) {
    const [selectedItem, setSelectedItem] = useState<LibraryItem | null>(null);
    const [isModalOpen, setIsModalOpen] = useState(false);

    const handleCardClick = (item: LibraryItem) => {
        setSelectedItem(item);
        setIsModalOpen(true);
    };

    return (
        <div className="mb-10">
            <h2 className="text-2xl font-bold text-[#e5a00d] mb-6">{name}</h2>

            {items && items.length > 0 ? (
                <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
                    {items.map((item, index) => (
                        <div
                            key={`${item.title}-${index}`}
                            className="bg-[#2d2d2d] rounded-lg overflow-hidden cursor-pointer transform transition-transform duration-200 hover:-translate-y-1"
                            onClick={() => handleCardClick(item)}
                        >
                            <div className="aspect-[2/3] relative">
                                {item.thumb_url ? (
                                    <img
                                        src={item.thumb_url}
                                        alt={item.title}
                                        className="w-full h-full object-cover"
                                        loading="lazy"
                                    />
                                ) : (
                                    <div className="w-full h-full bg-[#1a1a1a]" />
                                )}
                            </div>

                            <div className="p-3">
                                <h3 className="text-white text-sm font-medium line-clamp-2">
                                    {item.title}
                                </h3>
                                {item.year && (
                                    <p className="text-[#b3b3b3] text-sm mt-1">{item.year}</p>
                                )}
                            </div>
                        </div>
                    ))}
                </div>
            ) : (
                <div className="text-center p-6 bg-[#2d2d2d] rounded-lg text-[#b3b3b3]">
                    No items found in this library
                </div>
            )}

            <ItemModal
                item={selectedItem}
                isOpen={isModalOpen}
                onClose={() => {
                    setIsModalOpen(false);
                    setSelectedItem(null);
                }}
            />
        </div>
    );
}