export type Product = {
    id: string;
    title: string;
    description: string;
    duration: number;
    price: string;
    tags: string[];
};

export type CategoryWithProducts = {
    id: string;
    name: string;
    description: string;
    status: string;
    products: Product[];
};
