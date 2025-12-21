#!/bin/bash

# Crypto Tower Defense - Localhost Testing Environment
# This script sets up a complete testing environment with local blockchain

set -e

echo "🚀 Setting up Crypto Tower Defense Localhost Testing Environment"
echo "================================================================"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. Check if Anvil is running
echo -e "\n${BLUE}Step 1: Checking for running Anvil instance...${NC}"
if lsof -i:8545 > /dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  Port 8545 is already in use. Killing existing process...${NC}"
    lsof -ti:8545 | xargs kill -9 2>/dev/null || true
    sleep 2
fi

# 2. Start Anvil in background
echo -e "\n${BLUE}Step 2: Starting Anvil (local blockchain)...${NC}"
anvil --port 8545 --accounts 10 --balance 10000 > /tmp/anvil.log 2>&1 &
ANVIL_PID=$!
echo -e "${GREEN}✅ Anvil started (PID: $ANVIL_PID)${NC}"
sleep 3

# 3. Set default private key (Anvil account #0)
echo -e "\n${BLUE}Step 3: Setting up deployment account...${NC}"
export PRIVATE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
DEPLOYER_ADDRESS="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
echo -e "${GREEN}✅ Using Anvil account #0: $DEPLOYER_ADDRESS${NC}"

# 4. Deploy contracts
echo -e "\n${BLUE}Step 4: Deploying smart contracts...${NC}"
cd smart-contracts
forge script script/Deploy.s.sol --rpc-url http://localhost:8545 --broadcast

# 5. Load deployment addresses
echo -e "\n${BLUE}Step 5: Loading deployment addresses...${NC}"
if [ -f ".deployments.env" ]; then
    source .deployments.env
    echo -e "${GREEN}✅ Contracts deployed:${NC}"
    echo "  - GameToken: $GAME_TOKEN"
    echo "  - TowerToken: $TOWER_TOKEN"
    echo "  - CharacterNFT: $CHARACTER_NFT"
    echo "  - ItemNFT: $ITEM_NFT"
else
    echo -e "${YELLOW}⚠️  Deployment addresses file not found${NC}"
fi

# 6. Check backend status
echo -e "\n${BLUE}Step 6: Checking backend status...${NC}"
cd ..
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Backend is running on http://localhost:8080${NC}"
else
    echo -e "${YELLOW}⚠️  Backend not running. Start it with: ./start-backend.sh${NC}"
fi

# 7. Summary
echo -e "\n${GREEN}================================================================${NC}"
echo -e "${GREEN}✅ Localhost Testing Environment Ready!${NC}"
echo -e "${GREEN}================================================================${NC}"
echo ""
echo "🌐 Blockchain: http://localhost:8545 (Anvil)"
echo "🔧 Backend API: http://localhost:8080"
echo ""
echo "📝 Test Accounts (with 10000 ETH each):"
echo "  Account #0: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
echo "  Private Key: 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
echo ""
echo "📜 Deployed Contracts:"
echo "  - GameToken: $GAME_TOKEN"
echo "  - TowerToken: $TOWER_TOKEN"  
echo "  - CharacterNFT: $CHARACTER_NFT"
echo "  - ItemNFT: $ITEM_NFT"
echo ""
echo "🧪 Next Steps:"
echo "  1. Test contract interactions: cd smart-contracts && forge test -vvv"
echo "  2. Call contracts via cast: cast call $GAME_TOKEN 'name()' --rpc-url http://localhost:8545"
echo "  3. Test backend endpoints: curl http://localhost:8080/api/v1/auth/nonce"
echo ""
echo "🛑 To stop Anvil: kill $ANVIL_PID"
echo "📊 Anvil logs: tail -f /tmp/anvil.log"
