import type { PageServerLoad } from "./$types"
import type { Freelancer, Prestation, EventPickerMonth } from '$lib/types';
import { prestataires, prestations, monthsData } from '$lib/data';

type PageRes = {
    prestataires: Freelancer[];
    prestations: Prestation[];
    monthsData: EventPickerMonth[];
};

// TODO: change that function to use a fetch when the server is set
// async function getPrestataires(): Promise<Freelancer[]> {
//     return prestataires
// }
// NOTE: with the promise and shit
// export const load = async (): Promise<PageRes> => {
// export async function load(): Promise<PageRes> {
//     const prestataires = wawait getPrestataires()
//     return { prestataires, prestations, monthsData }
//
// }

// NOTE: the simple with mock data
export const load: PageServerLoad = (): PageRes => {
    return { prestataires, prestations, monthsData };
}
