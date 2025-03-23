type Photo = string;
// type Video = string
type Video = {
    thumbnail: string;
    duration: string;
};
export type EventPhotos = {
    date: string;
    photos: Photo[];
};
export type EventVideos = {
    date: string;
    videos: Video[];
};
