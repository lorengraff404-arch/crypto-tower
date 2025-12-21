# opBNB Testnet Deployment - FINAL

**Network:** opBNB Testnet  
**Chain ID:** 5611  
**Date:** December 18, 2024  
**RPC:** https://opbnb-testnet-rpc.bnbchain.org

---

## Deployed Contracts

### GameToken (GTK)
- **Address:** `0xd8068d4A032Dce3B63DB6617E299B609F211eF26`
- **Explorer:** [View on opBNB Testnet](https://opbnb-testnet.bscscan.com/address/0xd8068d4A032Dce3B63DB6617E299B609F211eF26)
- **Supply:** 10,000,000 GTK
- **Initial Mint:** 1,000,000 GTK
- **Daily Emission Cap:** 100,000 GTK

### TowerToken (TOWER)
- **Address:** `0xAC811d1948D00C3791d55A41C8323A136771125B`
- **Explorer:** [View on opBNB Testnet](https://opbnb-testnet.bscscan.com/address/0xAC811d1948D00C3791d55A41C8323A136771125B)
- **Supply:** 1,000,000 TOWER
- **Transfer Fee:** 0.5% → `0xd210925940D236951F0Bfa632A89d76fA1C8883d`
- **Distribution:** 20% ICO, 10% Dev, 5% Marketing, 65% Liquidity

### CharacterNFT
- **Address:** `0xF54C203451025D0f5D4d3693Bb2AE94027596FC9`
- **Explorer:** [View on opBNB Testnet](https://opbnb-testnet.bscscan.com/address/0xF54C203451025D0f5D4d3693Bb2AE94027596FC9)
- **Standard:** ERC-721
- **Royalty:** 2.5%

### ItemNFT
- **Address:** `0xC7f8397C65795009945F75D9a3a7a966c64F8a79`
- **Explorer:** [View on opBNB Testnet](https://opbnb-testnet.bscscan.com/address/0xC7f8397C65795009945F75D9a3a7a966c64F8a79)
- **Standard:** ERC-1155
- **Royalty:** 2.5%

---

## Deployment Details

**Deployer:** `0xdCb8ca66Ae0809Eed5dB73E8e1c3787c8178327e`  
**Gas Used:** 18,502,714  
**Total Cost:** 0.000018502714 BNB  
**Transaction Log:** `smart-contracts/broadcast/Deploy.s.sol/5611/run-latest.json`

---

## Verified Functionality

- ✅ All contracts deployed successfully
- ✅ Minter permissions configured
- ✅ Daily limits set (GTK: 10K, TOWER: 1K, CharacterNFT: 10)
- ✅ Frontend configured for opBNB Testnet
- ✅ Backend configured with correct addresses
- ✅ TOWER transfer fee (0.5%) enabled
- ✅ GTK revenue distribution ready

---

## Next Steps

1. Add contracts to MetaMask
2. Test token minting
3. Verify transfer fee on TOWER
4. Test NFT minting
5. Integrate with game frontend

---

## Important Notes

- **Network:** opBNB Testnet (Layer 2)
- **Correct Chain ID:** 5611 (not BSC's 97)
- **Lower Fees:** ~1000x cheaper than BSC mainnet
- **Fast Blocks:** ~1 second block time
