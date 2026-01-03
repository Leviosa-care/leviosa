import { type } from 'arktype'
import type { Infer } from 'sveltekit-superforms'

// Company Settings Schemas
export const companyNameSchema = type({
    name: "string"
})

export const companyEmailSchema = type({
    email: "string"
})

export const companyPhoneSchema = type({
    telephone: "string"
})

export const companyAddressSchema = type({
    address: "string"
})

export const companyInstagramSchema = type({
    instagram: "string"
})

// OTP Settings Schemas
export const otpDurationSchema = type({
    duration: "number"
})

export const otpLengthSchema = type({
    length: "number"
})

export const otpMaxAttemptsSchema = type({
    max_attempts: "number"
})

// Token Settings Schemas
export const accessTokenDurationSchema = type({
    duration: "number"
})

export const refreshTokenDurationSchema = type({
    duration: "number"
})

// Types
export type CompanyName = Infer<typeof companyNameSchema>
export type CompanyEmail = Infer<typeof companyEmailSchema>
export type CompanyPhone = Infer<typeof companyPhoneSchema>
export type CompanyAddress = Infer<typeof companyAddressSchema>
export type CompanyInstagram = Infer<typeof companyInstagramSchema>
export type OtpDuration = Infer<typeof otpDurationSchema>
export type OtpLength = Infer<typeof otpLengthSchema>
export type OtpMaxAttempts = Infer<typeof otpMaxAttemptsSchema>
export type AccessTokenDuration = Infer<typeof accessTokenDurationSchema>
export type RefreshTokenDuration = Infer<typeof refreshTokenDurationSchema>

// Defaults
export const companyNameDefaults: CompanyName = {
    name: "",
}

export const companyEmailDefaults: CompanyEmail = {
    email: "",
}

export const companyPhoneDefaults: CompanyPhone = {
    telephone: "",
}

export const companyAddressDefaults: CompanyAddress = {
    address: "",
}

export const companyInstagramDefaults: CompanyInstagram = {
    instagram: "",
}

export const otpDurationDefaults: OtpDuration = {
    duration: 300,
}

export const otpLengthDefaults: OtpLength = {
    length: 6,
}

export const otpMaxAttemptsDefaults: OtpMaxAttempts = {
    max_attempts: 5,
}

export const accessTokenDurationDefaults: AccessTokenDuration = {
    duration: 15,
}

export const refreshTokenDurationDefaults: RefreshTokenDuration = {
    duration: 168,
}
