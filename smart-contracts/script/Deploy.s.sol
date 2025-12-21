// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import {Script, console} from "forge-std/Script.sol";
import {GameToken} from "../src/GameToken.sol";
import {TowerToken} from "../src/TowerToken.sol";
import {CharacterNFT} from "../src/CharacterNFT.sol";
import {ItemNFT} from "../src/ItemNFT.sol";

contract DeployScript is Script {
    function run() external {
        // Load deployer private key from env
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);
        
        console.log("Deploying contracts with:", deployer);
        console.log("Deployer balance:", deployer.balance);

        vm.startBroadcast(deployerPrivateKey);

        // Deploy GameToken (GTK)
        console.log("\n=== Deploying GameToken ===");
        GameToken gtk = new GameToken();
        console.log("GameToken deployed to:", address(gtk));
        console.log("GTK Total Supply:", gtk.totalSupply() / 1e18, "tokens");

        // Deploy TowerToken (TOWER)
        console.log("\n=== Deploying TowerToken ===");
        address treasury = deployer;
        address buyback = deployer;
        address icoWallet = deployer;
        address devWallet = deployer;
        address marketingWallet = deployer;
        address liquidityWallet = deployer;
        
        TowerToken tower = new TowerToken(
            treasury,
            buyback,
            icoWallet,
            devWallet,
            marketingWallet,
            liquidityWallet
        );
        console.log("TowerToken deployed to:", address(tower));
        console.log("TOWER Total Supply:", tower.totalSupply() / 1e18, "tokens");

        // Deploy CharacterNFT
        console.log("\n=== Deploying CharacterNFT ===");
        string memory baseURI = "https://api.cryptotd.io/metadata/characters/";
        CharacterNFT characterNFT = new CharacterNFT(baseURI);
        console.log("CharacterNFT deployed to:", address(characterNFT));

        // Deploy ItemNFT
        console.log("\n=== Deploying ItemNFT ===");
        string memory itemURI = "https://api.cryptotd.io/metadata/items/{id}.json";
        ItemNFT itemNFT = new ItemNFT(itemURI);
        console.log("ItemNFT deployed to:", address(itemNFT));

        // Authorize deployer as minter with daily limits
        console.log("\n=== Setting up permissions ===");
        gtk.authorizeMinter(deployer, 10_000 * 10**18); // 10K GTK/day
        console.log("GTK: Authorized deployer as minter (10K/day)");
        
        tower.authorizeMinter(deployer, 1_000 * 10**18); // 1K TOWER/day
        console.log("TOWER: Authorized deployer as minter (1K/day)");
        
        characterNFT.authorizeMinter(deployer, 10); // 10 NFTs/day
        console.log("CharacterNFT: Authorized deployer as minter (10/day)");
        
        itemNFT.authorizeMinter(deployer);
        console.log("ItemNFT: Authorized deployer as minter");

        vm.stopBroadcast();

        // Log deployment addresses
        console.log("\n=== Deployment Complete ===");
        console.log("GameToken:", address(gtk));
        console.log("TowerToken:", address(tower));
        console.log("CharacterNFT:", address(characterNFT));
        console.log("ItemNFT:", address(itemNFT));
        console.log("\nDeployment addresses saved to .deployments.env");
    }
}
