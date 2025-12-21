// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../src/TokenSwap.sol";
import "forge-std/Script.sol";

contract DeploySwap is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        
        vm.startBroadcast(deployerPrivateKey);
        
        // Get existing token addresses from environment or hardcode
        address towerAddress = 0x66c83Ff8839D4DAD387AD22D48d48605D5C1aDa3;
        address gtkAddress = 0xDc4Bf3BbADC42Cc9feE3790B84ca1D532C036588;
        
        // Deploy swap contract
        TokenSwap swap = new TokenSwap(towerAddress, gtkAddress);
        
        console.log("TokenSwap deployed to:", address(swap));
        console.log("TOWER Token:", towerAddress);
        console.log("GTK Token:", gtkAddress);
        console.log("Ratio: 1 TOWER = 100 GTK");
        
        vm.stopBroadcast();
    }
}
