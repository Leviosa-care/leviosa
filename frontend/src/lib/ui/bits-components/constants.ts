// Button variants and sizes for type safety and maintainability
export const BUTTON_VARIANTS = {
    DEFAULT: "default",
    GHOST: "ghost",
    OUTLINE: "outline",
    CHEVRON: "chevron",
    ARROW: "arrow"
} as const;

export const BUTTON_SIZES = {
    DEFAULT: "default",
    SM: "sm",
    LG: "lg",
    ICON: "icon"
} as const;

// Card variants and sizes
export const CARD_VARIANTS = {
    DEFAULT: "default",
    ELEVATED: "elevated",
    OUTLINE: "outline"
} as const;

export const CARD_SIZES = {
    DEFAULT: "default",
    SM: "sm",
    LG: "lg"
} as const;

// Input variants and sizes
export const INPUT_VARIANTS = {
    DEFAULT: "default",
    ERROR: "error"
} as const;

export const INPUT_SIZES = {
    DEFAULT: "default",
    SM: "sm",
    LG: "lg"
} as const;

// Export TypeScript types for autocompletion and type checking
export type ButtonVariantType = typeof BUTTON_VARIANTS[keyof typeof BUTTON_VARIANTS];
export type ButtonSizeType = typeof BUTTON_SIZES[keyof typeof BUTTON_SIZES];
export type CardVariantType = typeof CARD_VARIANTS[keyof typeof CARD_VARIANTS];
export type CardSizeType = typeof CARD_SIZES[keyof typeof CARD_SIZES];
export type InputVariantType = typeof INPUT_VARIANTS[keyof typeof INPUT_VARIANTS];
export type InputSizeType = typeof INPUT_SIZES[keyof typeof INPUT_SIZES];
