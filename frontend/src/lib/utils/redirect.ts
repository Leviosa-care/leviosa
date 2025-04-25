export const MESSAGES = {
    basic: "Tu dois être connecté pour accéder à cette page",
    expiredSession: "Votre session a expiré. Veuillez vous reconnecter pour continuer",
} as const;

export type MessageType = keyof typeof MESSAGES
export const handleLoginRedirect = (url: URL, message: MessageType = "basic") => {
    const redirectTo = url.pathname + url.searchParams
    return `/auth?redirectTo=${redirectTo}&message=${message}`
}
