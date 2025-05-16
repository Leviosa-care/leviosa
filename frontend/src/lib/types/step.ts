

type Status = "complete" | "current" | "upcoming";
export type Step = {
    id: number;
    name: string;
    status: Status;
};
