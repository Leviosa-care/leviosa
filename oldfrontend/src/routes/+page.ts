// TODO: put that in the type folder to make that thing cleann brother
type FormControlType = 'email' | 'password';
type FormControl = {
    name: string;
    label: string;
    type: FormControlType;
    placeholder: string;
    value?: string;
};

const SignInFormControls: FormControl[] = [
    {
        name: 'email',
        label: 'Email',
        type: 'email',
        placeholder: 'Entre ton adresse email',
        // TODO: change that value when testing the backend
        value: 'admin@gmail.fr'
        // value: '',
    },
    {
        name: 'password',
        label: 'Mot de passe',
        type: 'password',
        placeholder: 'Entre ton mot de passe'
    }
];
const SignUpFormControl: FormControl[] = [
    {
        name: 'email',
        label: 'Email',
        type: 'email',
        placeholder: 'Entre ton adresse email',
        // TODO: change that value when testing the backend
        value: 'admin@gmail.fr'
        // value: '',
    }
];

type PageRes = { SignInFormControls: FormControl[]; SignUpFormControl: FormControl[] };
export function load(): PageRes {
    return { SignInFormControls, SignUpFormControl };
}
