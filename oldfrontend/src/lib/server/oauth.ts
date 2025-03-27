import { Google } from 'arctic';
// import { GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET } from '$env/static/private';

// export const google = new Google(
//     GOOGLE_CLIENT_ID,
//     GOOGLE_CLIENT_SECRET,
//     'http://localhost:5173/oauth/google/callback'
// );

console.log("from the oauth thing:", import.meta.env.VITE_GOOGLE_CLIENT_ID)
export const google = new Google(
    import.meta.env.VITE_GOOGLE_CLIENT_ID,
    import.meta.env.VITE_GOOGLE_CLIENT_SECRET,
    'http://localhost:5173/oauth/google/callback'
);
