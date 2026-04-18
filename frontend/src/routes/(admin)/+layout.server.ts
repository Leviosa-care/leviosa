export const load = async ({ locals, parent }) => {
    const { user, permissions } = await parent();

    return {
        user,
        role: user.role,
        permissions,
    };
};
