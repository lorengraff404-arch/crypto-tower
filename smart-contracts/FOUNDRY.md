# Foundry Setup

## Instalación Completada ✅

Foundry es una herramienta de desarrollo Ethereum escrita en Rust:
- **Más rápida** que Hardhat (10-100x)
- **Más estable** (menos conflictos de dependencias)
- **Testing** en Solidity nativo
- **Gas reports** integrados

## Comandos Foundry

### Compilar
```bash
forge build
```

### Test
```bash
forge test
forge test -vvv  # Verbose output
```

### Deploy
```bash
# Local
forge create --rpc-url http://localhost:8545 \
  --private-key $PRIVATE_KEY \
  src/GameToken.sol:GameToken

# BSC Testnet  
forge create --rpc-url https://data-seed-prebsc-1-s1.binance.org:8545 \
  --private-key $PRIVATE_KEY \
  --verify \
  src/GameToken.sol:GameToken
```

### Gas Report
```bash
forge test --gas-report
```

## Estructura del Proyecto

```
smart-contracts/
├── src/               # Contratos Solidity
│   ├── GameToken.sol
│   ├── TowerToken.sol
│   ├── CharacterNFT.sol
│   └── ItemNFT.sol
├── test/              # Tests en Solidity
├── script/            # Deploy scripts
├── out/               # Artifacts compilados
├── lib/               # Dependencies (OpenZeppelin)
└── foundry.toml       # Config
```

## Instalar OpenZeppelin

```bash
forge install OpenZeppelin/openzeppelin-contracts@v5.4.0
```

Luego actualizar remappings en `foundry.toml`:
```toml
remappings = [
    "@openzeppelin/=lib/openzeppelin-contracts/"
]
```

## Ventajas vs Hardhat

1. **Compilación**: ~100x más rápido
2. **Testing**: Escrito en Solidity (no JS)
3. **Dependencies**: Sin npm hell
4. **Gas**: Reportes precisos
5. **Debugging**: Stack traces mejores

## Migración de Hardhat

Ya movimos los contratos:
- `contracts/*.sol` → `src/*.sol`
- Hardhat artifacts → Foundry `out/`
- `hardhat.config.js` → `foundry.toml`

Hardhat aún disponible para scripts TypeScript si necesario.
