import { setContext, getContext } from "svelte";

export type Gender = '' | 'man' | 'woman' | 'non_binary' | 'prefer_not_to_say' | 'custom' | 'not precised';

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

type RegisterState = {
    email: string;
    password: string;
    address1: string;
    address2: string;
    telephone: string;
    postalCode: string;
    city: string;
    firstname: string;
    lastname: string;
    gender: Gender;
    birthdate: string;
};

export class Register {
    value = $state<RegisterState>({
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

    constructor(initial?: Partial<RegisterState>) {
        if (initial) {
            Object.assign(this.value, initial);
        }
    }

    static from(json: string | object): Register {
        const data = typeof json === "string" ? JSON.parse(json) : json;
        return new Register(data);
    }
    toJSON(): RegisterState {
        return { ...this.value };
    }
    persist() {
        sessionStorage.setItem("registration", JSON.stringify(this.toJSON()));
    }
    static load(): Register {
        const json = sessionStorage.getItem("registration");
        if (!json) return new Register();
        return Register.from(json);
    }
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
    all(): RegisterState {
        return this.value
    }
    clear(): void {
        sessionStorage.removeItem("registration");
        Object.assign(this.value, {
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
        });
    }
}

const REGISTRATION_KEY = Symbol('registration')

export function setRegistrationContext(): Register {
    const instance = Register.load();
    setContext(REGISTRATION_KEY, instance);
    return instance;
}

export function getRegistrationContext() {
    return getContext<ReturnType<typeof setRegistrationContext>>(REGISTRATION_KEY)
}

