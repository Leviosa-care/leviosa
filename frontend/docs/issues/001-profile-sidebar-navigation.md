# 001 — Profile sidebar navigation

**Type:** AFK
**Status:** open
**Blocked by:** None

## What to build

Add a "Mon profil" navigation entry to the staff sidebar so partners can reach `/staff/profile` from anywhere in the portal without guessing the URL.

The entry should appear in both the desktop collapsed/expanded sidebar and the mobile bottom navigation bar. Place it between the statistics entries and the settings icon. Use a user/person icon consistent with the existing Lucide icon set. The active-state highlight logic already used for other nav items should apply without modification.

## Acceptance criteria

- [ ] "Mon profil" link appears in the desktop sidebar navigation list
- [ ] "Mon profil" link appears in the mobile bottom navigation bar
- [ ] Link navigates to `/staff/profile`
- [ ] Active state is highlighted when the current path is `/staff/profile`
- [ ] Both `administrator` and `partner` roles see the link
- [ ] `npm run check` passes with no new type errors

## Parent

docs/prd/staff-pages-implementation.md
