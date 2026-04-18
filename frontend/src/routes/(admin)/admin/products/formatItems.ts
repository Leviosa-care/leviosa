export function formatItems(
    itemsList: Set<string>,
): { value: string; label: string }[] {
    const items: { value: string; label: string }[] = [];
    for (const item of itemsList) {
        items.push({ value: item, label: item });
    }
    return items;
}
