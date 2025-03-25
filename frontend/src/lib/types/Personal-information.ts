import type { Component } from 'svelte';

export type FieldValue = {
    name: string;
    value: string | number;
    properties?: any;
};

export type FieldConstructor = {
    name: string;
    fieldname: string;
    missingLabel?: string;
    addLabel?: string;
    modifyLabel: string;
    modifiedSlot: Component;
};

export type Field = FieldValue & FieldConstructor;
