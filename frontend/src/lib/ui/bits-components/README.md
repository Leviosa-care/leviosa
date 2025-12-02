# Bits UI Components

This directory contains refactored components built with `bits-ui` primitives and your design tokens from `src/app.css`.

## Components Overview

### Navigation & Layout
- **Tabs**: Tabbed navigation with list, triggers, and content panels
- **BackButton**: Navigation back button with customizable variants
- **Card**: Base card component with variants (default, elevated, outline)

### Display Components
- **EventCard**: Event display card with image, details, and metadata
- **ConsultationCard**: Service/consultation card with rating and pricing
- **NoEventCard**: Empty state card with call-to-action
- **CarouselIndicator**: Pagination indicators for carousels

### Form Components
- **FormInput**: Input field with label, validation, and variants
- **FormSelect**: Select dropdown with search and keyboard navigation
- **FormRadio**: Radio group with horizontal/vertical orientation
- **TimePickerButton**: Time selection button with icon

## Usage Examples

### Tabs Component
```svelte
<script>
    import { Tabs, TabsList, TabsTrigger, TabsContent } from '$lib/ui/bits-components';
    let selectedTab = 'events';
</script>

<Tabs bind:value={selectedTab}>
    <TabsList>
        <TabsTrigger value="events">Events</TabsTrigger>
        <TabsTrigger value="consultations">Consultations</TabsTrigger>
        <TabsTrigger value="profile">Profile</TabsTrigger>
    </TabsList>

    <TabsContent value="events">
        <!-- Events content -->
    </TabsContent>

    <TabsContent value="consultations">
        <!-- Consultations content -->
    </TabsContent>
</Tabs>
```

### EventCard Component
```svelte
<script>
    import { EventCard } from '$lib/ui/bits-components';

    function handleEventClick(id) {
        console.log('Event clicked:', id);
    }
</script>

<EventCard
    id="event-123"
    eventType="Workshop"
    title="Introduction to Meditation"
    date="Friday 21 June"
    eventBegin="10:00"
    city="Ivry-sur-Seine, Paris"
    prestationType="Meditation"
    prestationStart="11:00"
    onclick={() => handleEventClick('event-123')}
/>
```

### Form Components
```svelte
<script>
    import { FormInput, FormSelect, FormRadio } from '$lib/ui/bits-components';

    let email = '';
    let eventType = '';
    let preferences = '';

    const eventTypes = [
        { value: 'meditation', label: 'Meditation' },
        { value: 'yoga', label: 'Yoga' },
        { value: 'workshop', label: 'Workshop' }
    ];

    const preferenceOptions = [
        { value: 'morning', label: 'Morning sessions' },
        { value: 'evening', label: 'Evening sessions' }
    ];
</script>

<FormInput
    id="email"
    name="email"
    type="email"
    label="Email address"
    placeholder="Enter your email"
    bind:value={email}
    required
/>

<FormSelect
    id="eventType"
    name="eventType"
    label="Event type"
    bind:value={eventType}
    options={eventTypes}
    placeholder="Select an event type"
    required
/>

<FormRadio
    name="preferences"
    label="Session preferences"
    bind:value={preferences}
    options={preferenceOptions}
    orientation="horizontal"
/>
```

## Design System Integration

All components use your existing design tokens from `src/app.css`:
- **Colors**: `--foreground`, `--background`, `--accent`, `--destructive`, etc.
- **Spacing**: Consistent padding and margins
- **Typography**: Inter font family with proper sizing
- **Shadows**: `--shadow-card`, `--shadow-popover`, etc.
- **Borders**: `--border-card`, `--border-input`, etc.
- **Animations**: Smooth transitions and hover effects

## Accessibility

- **Keyboard navigation**: Full keyboard support for all interactive elements
- **ARIA attributes**: Proper labels and roles for screen readers
- **Focus management**: Visible focus indicators and logical tab order
- **Color contrast**: All colors meet WCAG AA standards

## Migration Path

To replace old components:

1. Import from `'$lib/ui/bits-components'` instead of old paths
2. Update props to match new component APIs
3. Test functionality and adjust as needed
4. Move components from `bits-components/` to main `ui/` directory when ready

## Features

- **Bits UI Primitives**: Built on top of robust bits-ui components
- **Design Tokens**: Consistent styling with your brand
- **TypeScript**: Full type safety and IntelliSense
- **Responsive**: Works on all screen sizes
- **Dark Mode**: Automatic theme switching
- **Animations**: Smooth transitions and interactions
- **Accessibility**: WCAG compliant components