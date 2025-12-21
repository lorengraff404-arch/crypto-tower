#!/bin/bash

# Script to add auth.js to all game HTML files
# This adds the auth script tag before the closing </body> tag if not already present

GAME_CLIENT_DIR="/Users/lorengraff/Development/1.-EN-DESARROLLO/Crypto_Towell_Defense/game-client"

# Files to update (excluding index.html which is the login page)
FILES=(
    "marketplace.html"
    "gacha.html"
)

AUTH_SCRIPT='    <!-- Auth System -->\n    <script src="js/auth.js"></script>'

for file in "${FILES[@]}"; do
    filepath="$GAME_CLIENT_DIR/$file"
    
    if [ -f "$filepath" ]; then
        # Check if auth.js is already included
        if grep -q "auth.js" "$filepath"; then
            echo "✓ $file already has auth.js"
        else
            # Find the last </body> tag and add auth.js before it
            # Using perl for in-place editing
            perl -i -pe 's{(</body>)}{    <!-- Auth System -->\n    <script src="js/auth.js"></script>\n\n$1}' "$filepath"
            echo "✓ Added auth.js to $file"
        fi
    else
        echo "✗ File not found: $file"
    fi
done

echo ""
echo "Done! Auth.js has been added to all game pages."
