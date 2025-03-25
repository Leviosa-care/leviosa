import { Writable } from 'svelte/store';

interface Persisted<T> extends Writable<T> {
    reset: () => void;
}
interface Serializer<T> {
    parse(text: string): T;
    stringify(object: T): string;
}
type StorageType = 'local' | 'session';
interface Options<StoreType, SerializerType> {
    serializer?: Serializer<SerializerType>;
    storage?: StorageType;
    syncTabs?: boolean;
    onError?: (e: unknown) => void;
    onWriteError?: (e: unknown) => void;
    onParseError?: (newValue: string | null, e: unknown) => void;
    beforeRead?: (val: SerializerType) => StoreType;
    beforeWrite?: (val: StoreType) => SerializerType;
}
/** @deprecated `writable()` has been renamed to `persisted()` */
declare function writable<StoreType, SerializerType = StoreType>(key: string, initialValue: StoreType, options?: Options<StoreType, SerializerType>): Persisted<StoreType>;
declare function persisted<StoreType, SerializerType = StoreType>(key: string, initialValue: StoreType, options?: Options<StoreType, SerializerType>): Persisted<StoreType>;

export { type Options, type Persisted, type Serializer, type StorageType, persisted, writable };
