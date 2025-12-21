#!/bin/bash

# Script to add modal system to all game HTML files

GAME_CLIENT_DIR="/Users/lorengraff/Development/1.-EN-DESARROLLO/Crypto_Towell_Defense/game-client"

# Files to update
FILES=(
    "index.html"
    "game.html"
    "island-raids.html"
    "marketplace.html"
)

echo "Adding modal system to all game pages..."
echo ""

for file in "${FILES[@]}"; do
    filepath="$GAME_CLIENT_DIR/$file"
    
    if [ -f "$filepath" ]; then
        # Check if modals.css is already included
        if grep -q "modals.css" "$filepath"; then
            echo "✓ $file already has modals.css"
        else
            # Add modals.css before </head>
            perl -i -pe 's{(</head>)}{    <link rel="stylesheet" href="css/modals.css">\n$1}' "$filepath"
            echo "✓ Added modals.css to $file"
        fi
        
        # Check if modals.js is already included
        if grep -q "modals.js" "$filepath"; then
            echo "✓ $file already has modals.js"
        else
            # Add modals.js before </body>
            perl -i -pe 's{(</body>)}{    <script src="js/modals.js"></script>\n$1}' "$filepath"
            echo "✓ Added modals.js to $file"
        fi
        
        echo ""
    else
        echo "✗ File not found: $file"
        echo ""
    fi
done

echo "Done! Modal system has been added to all game pages."
