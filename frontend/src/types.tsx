export interface LibraryItem {
    title: string;
    year?: number;
    thumb_url?: string;
    added_at: string;
    rating?: number;
    description?: string;
    cast?: string[];
    trailer_url?: string;
}

export interface Libraries {
    [key: string]: LibraryItem[];
}