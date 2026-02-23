# Petstore Frontend - Alpine.js + Tailwind CSS

## Overview

Client-side rendered (CSR) frontend for the Petstore REST API:
- **Alpine.js 3.x**: Reactive components, state management, all data fetching via `fetch()`
- **Tailwind CSS 3.x**: Utility-first styling via CDN
- No HTMX, no build step, no bundler

## Architecture

```
Browser
  ├── Static HTML/CSS/JS ← Go Fiber (app.Static("/", "../ui"))
  └── JSON data ← REST API endpoints (/pet, /category)
```

- **NOT SSR**: Go server only serves static files + REST API
- All rendering happens in the browser via Alpine.js directives
- API calls use `fetch()` returning JSON

## File Structure

```
ui/
├── agent.md           → This documentation
├── index.html         → Pet management page
└── categories.html    → Category management page
```

## Pages

### index.html - Pet Management
- **Alpine.js component**: `petStore()` with `x-data` and `x-init`
- **Features**: List, filter by status, create, edit, delete pets
- **Stats dashboard**: Total, available, pending, sold counts
- **Modals**: Add Pet, Edit Pet with form validation
- **API calls**: `fetch('/pet/findByStatus?status=...')`, `fetch('/pet', {method: 'POST', ...})`

### categories.html - Category Management
- **Alpine.js component**: `categoryStore()` with `x-data` and `x-init`
- **Features**: List, create, edit, delete categories
- **Inline editing**: Click edit → form replaces card content
- **API calls**: `fetch('/category/listAll')`, `fetch('/category', {method: 'POST', ...})`

## API Integration

### Endpoints Used

| Endpoint | Method | Used In |
|----------|--------|---------|
| `/pet/findByStatus?status=` | GET | index.html |
| `/pet` | POST | index.html (create) |
| `/pet` | PUT | index.html (update) |
| `/pet/:id` | DELETE | index.html |
| `/category/listAll` | GET | index.html + categories.html |
| `/category` | POST | categories.html |
| `/category` | PUT | categories.html |
| `/category/:id` | DELETE | categories.html |

### Error Handling
- Backend returns: `{"error": "message"}`
- UI parses: `error.error || error.message || 'Unexpected error'`
- States: loading spinner, error banner, empty state with CTA

### "All Pets" Filter
Makes 3 sequential `fetch()` calls (one per status: available, pending, sold) since there's no "list all" API endpoint.

## Alpine.js Pattern

```javascript
function petStore() {
    return {
        pets: [],
        loading: false,
        error: null,

        async init() {
            await this.loadPets();
        },

        async loadPets() {
            this.loading = true;
            const response = await fetch('/pet/findByStatus?status=available');
            if (response.ok) {
                this.pets = await response.json();
            }
            this.loading = false;
        }
    };
}
```

## Design System

- **Layout**: Gradient backgrounds, glassmorphism cards
- **Colors**: Blue-purple gradients (pets), purple-pink-red (categories)
- **Status badges**: Emerald (available), amber (pending), slate (sold)
- **Animations**: fadeIn 0.3s, slideDown 0.3s, bounce on icons
- **Responsive**: Mobile-first with Tailwind breakpoints (md, lg)

## Dependencies (CDN)

```html
<script src="https://cdn.tailwindcss.com"></script>
<script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
```

No build step, no node_modules, no bundler required.
