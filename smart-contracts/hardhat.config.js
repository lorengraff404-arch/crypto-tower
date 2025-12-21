import "@nomicfoundation/hardhat-ethers";
import "@nomicfoundation/hardhat-chai-matchers";
import "@nomicfoundation/hardhat-verify";
import "@typechain/hardhat";
import "hardhat-gas-reporter";
import "solidity-coverage";
// Note: @openzeppelin/hardhat-upgrades temporarily removed due to Node v24 compatibility
// Will add back when ecosystem fully supports Hardhat 3.x + Node v24
import dotenv from "dotenv";

dotenv.config();

/** @type import('hardhat/config').HardhatUserConfig */
const config = {
    solidity: {
        version: "0.8.28", // Latest stable Solidity
        settings: {
            optimizer: {
                enabled: true,
                runs: 200,
            },
            evmVersion: "cancun", // Latest EVM version
        },
    },
    networks: {
        // Local development
        hardhat: {
            chainId: 1337,
        },
        localhost: {
            url: "http://127.0.0.1:8545",
        },
        // opBNB Testnet
        opbnbTestnet: {
            url: "https://opbnb-testnet-rpc.bnbchain.org",
            chainId: 5611,
            accounts: process.env.PRIVATE_KEY ? [process.env.PRIVATE_KEY] : [],
            gasPrice: 1000000000, // 1 gwei
        },
        // BSC Testnet
        bscTestnet: {
            url: "https://data-seed-prebsc-1-s1.binance.org:8545",
            chainId: 97,
            accounts: process.env.PRIVATE_KEY ? [process.env.PRIVATE_KEY] : [],
            gasPrice: 10000000000, // 10 gwei
        },
        // opBNB Mainnet
        opbnb: {
            url: "https://opbnb-mainnet-rpc.bnbchain.org",
            chainId: 204,
            accounts: process.env.PRIVATE_KEY ? [process.env.PRIVATE_KEY] : [],
        },
        // BSC Mainnet
        bsc: {
            url: "https://bsc-dataseed.binance.org",
            chainId: 56,
            accounts: process.env.PRIVATE_KEY ? [process.env.PRIVATE_KEY] : [],
        },
    },
    etherscan: {
        apiKey: {
            bscTestnet: process.env.BSCSCAN_API_KEY || "",
            bsc: process.env.BSCSCAN_API_KEY || "",
            opbnbTestnet: process.env.BSCSCAN_API_KEY || "",
            opbnb: process.env.BSCSCAN_API_KEY || "",
        },
    },
    gasReporter: {
        enabled: process.env.REPORT_GAS === "true",
        currency: "USD",
    },
    typechain: {
        outDir: "typechain-types",
        target: "ethers-v6",
    },
    paths: {
        sources: "./contracts",
        tests: "./test",
        cache: "./cache",
        artifacts: "./artifacts",
    },
};

export default config;
