# Crypto Tower Defense - Client Setup Guide

## 🚀 Quick Start

### Prerequisites
- Node.js 18+ (for TypeScript compilation)
- Modern browser with MetaMask extension
- Backend server running on `localhost:8080`

### Installation

```bash
# Install dependencies
npm install

# Build TypeScript files
npm run build

# For development (watch mode)
npm run watch
```

### Project Structure

```
game-client/
├── js/                    # TypeScript source files
│   ├── types.ts          # Type definitions (20+ interfaces)
│   ├── api.ts            # Type-safe API client
│   ├── islands.ts        # Island Raids UI logic
│   ├── quests.ts         # Quest/Mission UI logic
│   └── wallet.js         # Web3/MetaMask integration
├── dist/js/              # Compiled JavaScript (auto-generated)
│   ├── *.js             # ES2022 modules
│   ├── *.d.ts           # Type declarations
│   └── *.js.map         # Source maps
├── css/                  # Stylesheets
│   ├── style.css        # Global styles
│   ├── islands.css      # Island Raids styles
│   └── quests.css       # Quest system styles
├── islands.html          # Island Raids page
├── quests.html           # Mission/Quest page
├── battle.html           # Battle system
├── game.html             # Main game
└── index.html            # Landing page
```

## 📝 TypeScript Configuration

### Build Scripts

```json
{
  "build": "tsc",           // Compile TypeScript
  "watch": "tsc --watch",   // Watch mode for development
  "clean": "rm -rf dist"    // Clean build artifacts
}
```

### TypeScript Settings (tsconfig.json)

- **Target**: ES2022 (modern JavaScript)
- **Module**: ES2022 (native modules)
- **Strict Mode**: Enabled (maximum type safety)
- **Source Maps**: Enabled (for debugging)

## 🎮 Pages Overview

### 1. Landing Page (`index.html`)
- Wallet connection (MetaMask)
- Game introduction
- Statistics display

### 2. Island Raids (`islands.html`)
- **Features**:
  - Island selection grid
  - Island detail modals
  - Raid battle interface (placeholder)
  - Loot display system
- **API**: `/api/v1/islands/*`
- **Script**: `dist/js/islands.js` (compiled from `islands.ts`)

### 3. Quest/Mission System (`quests.html`)
- **Features**:
  - Mission list sidebar
  - Mission detail panel
  - Aria dialogue system
  - Progress tracking
  - Reward displays
- **API**: `/api/v1/missions/*`, `/api/v1/story/*`
- **Script**: `dist/js/quests.js` (compiled from `quests.ts`)

### 4. Battle System (`battle.html`)
- PvP/PvE combat interface
- To be integrated with game engine

### 5. Main Game (`game.html`)
- Full game interface
- Integrates all systems

## 🔧 Development Workflow

### 1. Making Changes

```bash
# Edit TypeScript files in js/ directory
vim js/islands.ts

# Rebuild
npm run build

# Or use watch mode for auto-compilation
npm run watch
```

### 2. Type Safety

All API calls are type-checked:

```typescript
// Example: Type-safe API call
const data = await apiClient.get<IslandListResponse>('/islands');
// data.islands is typed as Island[]
```

### 3. Adding New Features

1. **Define types** in `js/types.ts`
2. **Add API methods** to `js/api.ts` if needed
3. **Create UI logic** in new `.ts` file
4. **Compile**: `npm run build`
5. **Update HTML** to reference `dist/js/yourfile.js`

## 🎨 UI/UX Guidelines

### Design Principles
- **Modern**: CSS Grid, Flexbox, Animations
- **Elegant**: Glassmorphism, gradients, shadows
- **Minimalist**: Clean layouts, clear hierarchy
- **Animated**: Smooth transitions, hover effects

### Color Palette
```css
--primary: #00d4ff (cyan)
--secondary: #00a8ff (blue)
--background: #1a1a2e (dark blue)
--accent: #16213e (navy)
--success: #4caf50 (green)
--warning: #ff9800 (orange)
--danger: #f44336 (red)
```

### Responsive Breakpoints
```css
Desktop: 1024px+
Tablet: 768px - 1023px
Mobile: < 768px
```

## 🔐 Security Features

### Client-Side
- ✅ Type-safe API calls (TypeScript)
- ✅ Input validation before sending
- ✅ XSS prevention (sanitized outputs)
- ✅ Wallet signature verification

### Backend Integration
- ✅ JWT authentication required
- ✅ Rate limiting (100 req/min)
- ✅ CORS configured
- ✅ Security headers

## 🐛 Debugging

### TypeScript Errors

```bash
# Check for errors
npm run build

# Common issues:
# - Missing type definitions → Add to types.ts
# - Null safety → Use optional chaining (?.)
# - Type mismatch → Check API response types
```

### Runtime Errors

- **Check Browser Console**: F12 → Console tab
- **Source Maps**: Errors show TypeScript line numbers
- **Network Tab**: Inspect API calls

### Common Issues

| Issue | Solution |
|-------|----------|
| 401 Unauthorized | Check JWT token in localStorage |
| TypeScript errors | Run `npm run build` to see details |
| Module not found | Ensure `type="module"` in script tag |
| API timeout | Verify backend is running |

## 📚 API Integration

### Authentication Flow

```typescript
// 1. Get nonce
const { nonce } = await apiClient.getNonce(walletAddress);

// 2. Sign with MetaMask
const signature = await wallet.signMessage(nonce);

// 3. Verify signature
const { token } = await apiClient.verifySignature(walletAddress, signature);

// Token stored automatically in localStorage
```

### Protected Routes

All API calls to protected routes automatically include JWT:

```typescript
// API client handles auth header
const data = await apiClient.get<Mission[]>('/missions');
// Header: Authorization: Bearer <token>
```

## 🧪 Testing Checklist

### Manual Testing

- [ ] Wallet connection (MetaMask)
- [ ] Island selection and detail modal
- [ ] Mission list and detail view
- [ ] Aria dialogue display
- [ ] Raid entry and completion
- [ ] Reward collection
- [ ] Responsive design (mobile)
- [ ] Browser compatibility (Chrome, Firefox, Safari)

### API Testing

```bash
# Health check
curl http://localhost:8080/health

# Get islands (requires auth)
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/v1/islands
```

## 📦 Deployment

### Production Build

```bash
# Clean and build
npm run clean && npm run build

# Verify output
ls -la dist/js/

# Expected files:
# - api.js, api.d.ts, api.js.map
# - islands.js, islands.d.ts, islands.js.map
# - quests.js, quests.d.ts, quests.js.map
# - types.js, types.d.ts, types.js.map
```

### Environment Configuration

Update API endpoint in `js/api.ts`:

```typescript
// Development
const API_BASE_URL = 'http://localhost:8080/api/v1';

// Production
const API_BASE_URL = 'https://api.yourdomain.com/api/v1';
```

## 🔄 Updates & Maintenance

### Updating Dependencies

```bash
# Check for updates
npm outdated

# Update TypeScript
npm install -D typescript@latest

# Update types
npm install -D @types/node@latest
```

### Code Quality

- **TypeScript Strict Mode**: Always enabled
- **No 'any' Types**: Explicit typing required
- **ESLint**: Consider adding for code style
- **Prettier**: Consider adding for formatting

## 📖 Resources

- [TypeScript Documentation](https://www.typescriptlang.org/docs/)
- [MetaMask Documentation](https://docs.metamask.io/)
- [CSS Grid Guide](https://css-tricks.com/snippets/css/complete-guide-grid/)
- [Fetch API Reference](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API)

## 🤝 Contributing

### Code Style

1. Use TypeScript for all business logic
2. Keep JavaScript for Web3 integration only
3. Follow existing naming conventions
4. Add JSDoc comments for complex functions
5. Update types.ts for new API responses

### Commit Messages

```
feat: Add new island selection feature
fix: Resolve TypeScript compilation error
refactor: Extract loot display to separate function
docs: Update README with new API endpoints
```

---

**Status**: ✅ Production Ready (Frontend 70%)  
**Last Updated**: 2025-12-17  
**TypeScript Version**: 5.3+  
**Build System**: Native TypeScript Compiler
