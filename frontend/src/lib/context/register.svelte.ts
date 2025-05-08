import { setContext, getContext } from "svelte";

export type Gender = 'man' | 'woman' | 'non_binary' | 'prefer_not_to_say' | 'custom' | 'not precised';

type Address = {
    address1: string;
    address2: string;
    postalCode: string;
    city: string;
}

type General = {
    firstname: string;
    lastname: string;
    gender: Gender;
    birthdate: string;
    telephone: string;
}

export class Register {
    value = $state({
        email: "",
        password: "",
        address1: "",
        address2: "",
        telephone: "",
        postalCode: "",
        city: "",
        firstname: "",
        lastname: "",
        gender: "prefer_not_to_say",
        birthdate: "",
    })
    constructor() { }
    setEmail(email: string) {
        this.value.email = email
    }
    setAddress(address: Address) {
        this.value.address1 = address.address1
        this.value.address2 = address.address2
        this.value.postalCode = address.postalCode
        this.value.city = address.city
    }
    setPassword(password: string) {
        this.value.password = password
    }
    setGeneral(general: General) {
        this.value.firstname = general.firstname;
        this.value.lastname = general.lastname;
        this.value.gender = general.gender as Gender;
        this.value.birthdate = general.birthdate;
        this.value.telephone = general.telephone
    }
    all() {
        return this.value
    }
}

const REGISTRATION_KEY = Symbol('registration')

export function setRegistrationContext(): Register {
    return setContext(REGISTRATION_KEY, new Register())
}

export function getRegistrationContext() {
    return getContext<ReturnType<typeof setRegistrationContext>>(REGISTRATION_KEY)
}
