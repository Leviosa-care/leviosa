// export function focusNextTarget(event: Event): void {
//     const input = event.target as HTMLInputElement
//     const numericValue = input.value.replace(/[^0-9]/g, '');
//     const currentID = parseInt(input.id)
//     const nextInput = document.getElementById(`${currentID + 1}`) as HTMLInputElement | null;
//
//     if (numericValue.length > 0) {
//         input.value = numericValue[0]
//         input.dispatchEvent(new Event('input'));
//
//         if (nextInput) {
//             nextInput.focus();
//         }
//     } else {
//         // input.value = "";
//         input.value = "";
//         input.dispatchEvent(new Event('input'));
//     }
// }

export function focusNextTarget(event: Event): void {
    const input = event.target as HTMLInputElement;
    const currentID = parseInt(input.id);
    const nextInput = document.getElementById(`${currentID + 1}`) as HTMLInputElement | null;

    // Prevent non-numeric characters — but don't update input.value directly
    if (!/^\d$/.test(input.value)) {
        input.value = '';
        return;
    }

    // Move to next input if valid
    if (input.value && nextInput) {
        nextInput.focus();
    }
}

export function focusPrevTarget(event: KeyboardEvent): void {
    const key = event.key.toLowerCase();
    if (key === "backspace" || key === "delete") {
        const input = event.target as HTMLInputElement;
        const currentID = parseInt(input.id)
        const isEmpty = input.value.length === 0
        if (isEmpty) {
            // Find previous input if it exists
            const prevInput = document.getElementById(`${currentID - 1}`) as HTMLInputElement | null;
            if (prevInput) {
                prevInput.focus();
            }
        }
    }
}
