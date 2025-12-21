# 🎨 MODAL SYSTEM - GUÍA DE USO

## ✅ Sistema Implementado

Se ha creado un sistema centralizado de modales/popups con estilo **premium dark moderno** que se aplica automáticamente a todas las páginas del juego.

---

## 📁 Archivos Creados

1. **`css/modals.css`** - Estilos premium dark para modales
2. **`js/modals.js`** - Utilidades JavaScript para gestionar modales

---

## 🎯 Páginas Actualizadas

✅ `index.html`  
✅ `game.html`  
✅ `island-raids.html`  
✅ `marketplace.html`

**Todos los archivos ahora incluyen:**
```html
<link rel="stylesheet" href="css/modals.css">
<script src="js/modals.js"></script>
```

---

## 💡 CÓMO USAR

### 1. Mostrar Modal Simple

```javascript
showModal({
    title: 'Welcome!',
    content: '<p>This is a premium dark modal!</p>',
    buttons: [
        {
            text: 'Close',
            type: 'primary'
        }
    ]
});
```

### 2. Mostrar Confirmación

```javascript
showConfirm({
    title: 'Delete Character?',
    message: 'Are you sure you want to delete this character? This action cannot be undone.',
    type: 'danger',
    confirmText: 'Delete',
    cancelText: 'Cancel',
    onConfirm: () => {
        // Delete character
        console.log('Character deleted');
    },
    onCancel: () => {
        console.log('Cancelled');
    }
});
```

### 3. Mostrar Notificación (Toast)

```javascript
// Success
showNotification('Character created successfully!', 'success');

// Error
showNotification('Failed to connect wallet', 'error');

// Warning
showNotification('Low balance warning', 'warning');

// Info
showNotification('New feature available', 'info');
```

### 4. Mostrar Alerta

```javascript
showAlert('Success', 'Your transaction was completed!', 'success');
showAlert('Error', 'Something went wrong', 'error');
showAlert('Warning', 'Please verify your wallet', 'warning');
```

### 5. Modal con Múltiples Botones

```javascript
showModal({
    title: 'Choose Action',
    content: '<p>What would you like to do?</p>',
    buttons: [
        {
            text: 'Cancel',
            type: 'secondary'
        },
        {
            text: 'Save Draft',
            type: 'primary',
            onClick: () => console.log('Saved as draft')
        },
        {
            text: 'Publish',
            type: 'success',
            onClick: () => console.log('Published')
        }
    ]
});
```

---

## 🎨 Tipos de Botones

- **`primary`** - Azul (acción principal)
- **`secondary`** - Gris (acción secundaria)
- **`success`** - Verde (confirmación positiva)
- **`danger`** - Rojo (acción destructiva)

---

## 🎨 Tipos de Notificaciones

- **`success`** ✓ - Verde con borde verde
- **`error`** ✗ - Rojo con borde rojo
- **`warning`** ⚠ - Naranja con borde naranja
- **`info`** ℹ - Azul con borde azul

---

## 🔧 Características

### Estilos Premium Dark:
- ✅ Fondo degradado oscuro (#1e293b → #0f172a)
- ✅ Bordes con glow azul (#3b82f6)
- ✅ Sombras profundas con blur
- ✅ Animaciones suaves (fade in, slide up)
- ✅ Backdrop blur en overlay

### Funcionalidades:
- ✅ Cerrar con botón X
- ✅ Cerrar con tecla ESC
- ✅ Cerrar clickeando fuera del modal
- ✅ Scroll interno si contenido es largo
- ✅ Responsive (mobile-friendly)
- ✅ Auto-cierre de notificaciones (4 segundos)

---

## 📱 Responsive

El sistema es completamente responsive:
- En móvil: modales ocupan 95% del ancho
- Botones se apilan verticalmente
- Notificaciones se adaptan al ancho de pantalla

---

## 🎯 Reemplazar Modales Antiguos

### Antes (estilo antiguo):
```javascript
alert('Hello!'); // ❌ Feo, no personalizable
```

### Ahora (estilo premium):
```javascript
showAlert('Welcome', 'Hello!', 'info'); // ✅ Premium dark
```

### Antes (confirm antiguo):
```javascript
if (confirm('Delete?')) { // ❌ Feo
    deleteItem();
}
```

### Ahora (confirm premium):
```javascript
showConfirm({
    title: 'Confirm Delete',
    message: 'Are you sure?',
    type: 'danger',
    onConfirm: () => deleteItem()
}); // ✅ Premium dark
```

---

## 🚀 Ejemplos Prácticos

### Login Success:
```javascript
showNotification('Wallet connected successfully!', 'success');
```

### Purchase Confirmation:
```javascript
showConfirm({
    title: 'Purchase Character',
    message: `Buy this SSS character for 1000 GTK?`,
    type: 'warning',
    confirmText: 'Buy Now',
    onConfirm: async () => {
        await purchaseCharacter();
        showNotification('Character purchased!', 'success');
    }
});
```

### Error Handling:
```javascript
try {
    await someAction();
} catch (error) {
    showAlert('Error', error.message, 'error');
}
```

### Info Modal:
```javascript
showModal({
    title: 'How to Play',
    content: `
        <h3>Game Rules:</h3>
        <ul>
            <li>Build your team</li>
            <li>Battle other players</li>
            <li>Earn rewards</li>
        </ul>
    `,
    buttons: [
        { text: 'Got it!', type: 'primary' }
    ]
});
```

---

## ✅ Ventajas del Sistema

1. **Consistencia** - Todos los modales tienen el mismo estilo premium
2. **Fácil de usar** - Funciones simples y claras
3. **Personalizable** - Múltiples opciones de configuración
4. **Responsive** - Funciona en todos los dispositivos
5. **Accesible** - Soporte para teclado (ESC para cerrar)
6. **Moderno** - Animaciones y efectos visuales premium

---

## 🎨 Personalización Avanzada

### Modal con HTML Personalizado:
```javascript
showModal({
    title: 'Character Stats',
    content: `
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 10px;">
            <div>HP: <strong>1000</strong></div>
            <div>ATK: <strong>250</strong></div>
            <div>DEF: <strong>150</strong></div>
            <div>SPD: <strong>80</strong></div>
        </div>
    `
});
```

### Notificación con Duración Personalizada:
```javascript
showNotification('This will stay for 10 seconds', 'info', 10000);
```

---

## 🔄 Migración de Código Antiguo

Busca en tu código y reemplaza:

```javascript
// ❌ Antiguo
alert('Message');
// ✅ Nuevo
showAlert('Notice', 'Message', 'info');

// ❌ Antiguo
if (confirm('Sure?')) { ... }
// ✅ Nuevo
showConfirm({ message: 'Sure?', onConfirm: () => { ... } });

// ❌ Antiguo
console.log('Success!');
// ✅ Nuevo
showNotification('Success!', 'success');
```

---

## 🎉 ¡Listo para Usar!

El sistema está completamente implementado y listo para usar en todas las páginas del juego. Solo llama a las funciones y disfruta de los modales premium dark!

