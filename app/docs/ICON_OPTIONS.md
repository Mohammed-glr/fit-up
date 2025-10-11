# 🎨 Beautiful Icon Options for Tab Bar

Your tab bar now uses **Ionicons** and **Material Community Icons** - modern, beautiful icon libraries!

## Current Icons (Already Applied)

### 📊 Dashboard
- **Current**: `stats-chart` (Ionicons) - Analytics chart icon
- **Alternatives**:
  - `analytics` - Line chart with curves
  - `pie-chart` - Circular analytics  
  - `pulse` - Heart rate/activity monitor style

### 📅 Schema
- **Current**: `calendar` (Ionicons) - Clean calendar icon
- **Alternatives**:
  - `calendar-outline` - Outlined version
  - `today` (Material) - Calendar with today marker

### 💬 Messages  
- **Current**: `chatbubbles` (Ionicons) - Multiple chat bubbles
- **Alternatives**:
  - `chatbox-ellipses` - Single chat with typing indicator
  - `mail` - Email/message icon
  - `send` - Paper plane icon

### 🧠 Mindfulness
- **Current**: `brain` (Material Community) - Brain icon
- **Alternatives**:
  - `flower-outline` - Lotus flower (meditation symbol)
  - `spa` (Material) - Spa/wellness icon
  - `leaf-outline` - Nature/calm icon
  - `heart-pulse` (Material) - Heart with pulse line

### 👤 Profile
- **Current**: `person-circle` (Ionicons) - Person in circle
- **Alternatives**:
  - `person` - Simple person outline
  - `person-outline` - Outlined person

## How to Change Icons

Edit `app/(tabs)/_layout.tsx` and change the `name` prop in `IconSymbol`:

```tsx
<IconSymbol 
  size={focused ? 28 : 24} 
  name="analytics"  // Change this to any icon from the list
  color={color} 
/>
```

## Icon Libraries Used

1. **Ionicons** - Modern, clean iOS-style icons
2. **Material Community Icons** - Extended Material Design icons with 6000+ options
3. **FontAwesome 5** - Popular icon library with tons of options

## Browse More Icons

- Ionicons: https://ionic.io/ionicons
- Material Community: https://pictogrammers.com/library/mdi/
- FontAwesome: https://fontawesome.com/icons
- All in Expo: https://icons.expo.fyi/

## 🎯 Recommended Icon Combinations

### Option 1: Analytics & Wellness
- Dashboard: `analytics`
- Schema: `calendar`
- Messages: `chatbubbles`
- Mindfulness: `flower-outline`
- Profile: `person-circle`

### Option 2: Modern & Bold
- Dashboard: `pulse`
- Schema: `calendar`
- Messages: `chatbox-ellipses`
- Mindfulness: `spa`
- Profile: `person`

### Option 3: Fitness Focus
- Dashboard: `fitness`
- Schema: `calendar`
- Messages: `chatbubbles`
- Mindfulness: `heart-pulse`
- Profile: `person-circle`
