import type { PageServerLoad } from "./$types"
import { redirect } from '@sveltejs/kit';
import type { RequestEvent } from './$types';

export type QRCodeType = string;

// const name: string = "John"
const qrcode: QRCodeType =
    'https://1.bp.blogspot.com/-dHN4KiD3dsU/XRxU5JRV7DI/AAAAAAAAAz4/u1ynpCMIuKwZMA642dHEoXFVKuHQbJvwgCEwYBhgL/s1600/qr-code.png';

type PageRes = { name: string | undefined; qrcode: QRCodeType };
export const load: PageServerLoad = ({ locals }: RequestEvent): PageRes => {
    if (locals.user === null) {
        return redirect(302, '/');
    }
    return { name: locals.user?.firstname, qrcode };
}

// TODO: find a way to do return the async with the qrcode
// export async function load(): Promise<QRCode> {
//     const qrcode = await fetch("some api")
//     return { qrcode }
// }

// here to fetch something and to see what is going on brother
// fetch('https://jsonplaceholder.typicode.com/todos/1')
//     .then(response => response.json())
//     .then(json => console.log(json))
