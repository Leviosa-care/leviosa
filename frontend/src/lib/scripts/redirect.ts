import { goto } from '$app/navigation';

export function redirectTo(pathname: string) {
    goto(pathname);
}
