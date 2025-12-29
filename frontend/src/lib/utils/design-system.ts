import { tv, type VariantProps } from 'tailwind-variants';
import { twMerge } from 'tailwind-merge';

/**
 * Utility function to merge Tailwind classes with proper conflict resolution
 */
export function cn(...inputs: any[]) {
    return twMerge(inputs);
}

/**
 * Button component variants using tailwind-variants
 */
export const buttonVariants = tv({
    base: [
        'inline-flex items-center justify-center gap-2',
        'rounded-button px-4 py-2',
        'text-sm font-medium',
        'transition-colors duration-200',
        'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2',
        'disabled:pointer-events-none disabled:opacity-50'
    ],
    variants: {
        variant: {
            default: [
                'bg-dark text-contrast',
                'hover:bg-dark/90',
                'focus-visible:ring-dark'
            ],
            destructive: [
                'bg-destructive text-contrast',
                'hover:bg-destructive/90',
                'focus-visible:ring-destructive'
            ],
            outline: [
                'border border-input bg-background',
                'hover:bg-accent hover:text-accent-foreground',
                'focus-visible:ring-dark'
            ],
            secondary: [
                'bg-muted text-muted-foreground',
                'hover:bg-muted/80',
                'focus-visible:ring-muted'
            ],
            ghost: [
                'hover:bg-accent hover:text-accent-foreground',
                'focus-visible:ring-accent'
            ],
            link: [
                'text-dark underline-offset-4',
                'hover:underline',
                'focus-visible:ring-dark'
            ]
        },
        size: {
            default: 'h-10 px-4 py-2',
            sm: 'h-9 rounded-9px px-3',
            lg: 'h-11 rounded-10px px-8',
            icon: 'h-10 w-10'
        }
    },
    defaultVariants: {
        variant: 'default',
        size: 'default'
    }
});

/**
 * Alert component variants
 */
export const alertVariants = tv({
    base: [
        'relative w-full rounded-15px border p-4',
        'shadow-mini'
    ],
    variants: {
        variant: {
            default: [
                'bg-background text-foreground',
                'border-border'
            ],
            destructive: [
                'border-destructive/50 text-destructive',
                'dark:border-destructive'
            ],
            warning: [
                'border-tertiary/50 text-tertiary',
                'bg-tertiary/10'
            ],
            success: [
                'border-accent/50 text-accent-foreground',
                'bg-accent/10'
            ]
        }
    },
    defaultVariants: {
        variant: 'default'
    }
});

/**
 * Input component variants
 */
export const inputVariants = tv({
    base: [
        'flex w-full rounded-input border border-input',
        'bg-background px-3 py-2 text-sm',
        'ring-offset-background',
        'file:border-0 file:bg-transparent file:text-sm file:font-medium',
        'placeholder:text-muted-foreground',
        'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-dark focus-visible:ring-offset-2',
        'disabled:cursor-not-allowed disabled:opacity-50'
    ],
    variants: {
        variant: {
            default: 'hover:border-border-input-hover',
            error: 'border-destructive focus-visible:ring-destructive'
        },
        size: {
            default: 'h-10',
            sm: 'h-9',
            lg: 'h-11'
        }
    },
    defaultVariants: {
        variant: 'default',
        size: 'default'
    }
});

/**
 * Card component variants
 */
export const cardVariants = tv({
    base: [
        'rounded-card border border-card bg-background',
        'text-foreground shadow-card'
    ],
    variants: {
        variant: {
            default: '',
            elevated: 'shadow-popover',
            outline: 'border-border'
        },
        size: {
            default: 'p-6',
            sm: 'p-4',
            lg: 'p-8'
        }
    },
    defaultVariants: {
        variant: 'default',
        size: 'default'
    }
});

// Export types for component props
export type ButtonVariants = VariantProps<typeof buttonVariants>;
export type AlertVariants = VariantProps<typeof alertVariants>;
export type InputVariants = VariantProps<typeof inputVariants>;
export type CardVariants = VariantProps<typeof cardVariants>;