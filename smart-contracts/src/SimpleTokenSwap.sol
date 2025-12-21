// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IERC20 {
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
    function transfer(address to, uint256 amount) external returns (bool);
    function balanceOf(address account) external view returns (uint256);
}

contract SimpleTokenSwap {
    IERC20 public immutable towerToken;
    IERC20 public immutable gtkToken;
    address public owner;
    
    // Fixed ratio: 1 TOWER = 100 GTK
    uint256 public constant RATIO = 100;
    
    // Security: Rate limiting
    mapping(address => uint256) public lastSwapTime;
    uint256 public constant MIN_SWAP_INTERVAL = 60; // 1 minute
    
    // Security: Per-transaction limits
    uint256 public maxTowerPerSwap = 1000 * 10**18;
    uint256 public maxGTKPerSwap = 100000 * 10**18;
    
    event SwapTowerForGTK(address indexed user, uint256 towerAmount, uint256 gtkAmount);
    event SwapGTKForTower(address indexed user, uint256 gtkAmount, uint256 towerAmount);
    event LiquidityAdded(uint256 towerAmount, uint256 gtkAmount);
    
    modifier onlyOwner() {
        require(msg.sender == owner, "Not owner");
        _;
    }
    
    constructor(address _towerToken, address _gtkToken) {
        require(_towerToken != address(0) && _gtkToken != address(0), "Invalid addresses");
        towerToken = IERC20(_towerToken);
        gtkToken = IERC20(_gtkToken);
        owner = msg.sender;
    }
    
    function swapTowerForGTK(uint256 towerAmount) external {
        require(towerAmount > 0 && towerAmount <= maxTowerPerSwap, "Invalid amount");
        require(block.timestamp >= lastSwapTime[msg.sender] + MIN_SWAP_INTERVAL, "Too frequent");
        
        uint256 gtkAmount = towerAmount * RATIO;
        require(gtkToken.balanceOf(address(this)) >= gtkAmount, "Insufficient GTK liquidity");
        
        require(towerToken.transferFrom(msg.sender, address(this), towerAmount), "TOWER transfer failed");
        require(gtkToken.transfer(msg.sender, gtkAmount), "GTK transfer failed");
        
        lastSwapTime[msg.sender] = block.timestamp;
        emit SwapTowerForGTK(msg.sender, towerAmount, gtkAmount);
    }
    
    function swapGTKForTower(uint256 gtkAmount) external {
        require(gtkAmount > 0 && gtkAmount <= maxGTKPerSwap, "Invalid amount");
        require(gtkAmount % RATIO == 0, "GTK must be divisible by 100");
        require(block.timestamp >= lastSwapTime[msg.sender] + MIN_SWAP_INTERVAL, "Too frequent");
        
        uint256 towerAmount = gtkAmount / RATIO;
        require(towerToken.balanceOf(address(this)) >= towerAmount, "Insufficient TOWER liquidity");
        
        require(gtkToken.transferFrom(msg.sender, address(this), gtkAmount), "GTK transfer failed");
        require(towerToken.transfer(msg.sender, towerAmount), "TOWER transfer failed");
        
        lastSwapTime[msg.sender] = block.timestamp;
        emit SwapGTKForTower(msg.sender, gtkAmount, towerAmount);
    }
    
    function addLiquidity(uint256 towerAmount, uint256 gtkAmount) external onlyOwner {
        if (towerAmount > 0) {
            require(towerToken.transferFrom(msg.sender, address(this), towerAmount), "TOWER transfer failed");
        }
        if (gtkAmount > 0) {
            require(gtkToken.transferFrom(msg.sender, address(this), gtkAmount), "GTK transfer failed");
        }
        emit LiquidityAdded(towerAmount, gtkAmount);
    }
    
    function getReserves() external view returns (uint256 towerReserve, uint256 gtkReserve) {
        towerReserve = towerToken.balanceOf(address(this));
        gtkReserve = gtkToken.balanceOf(address(this));
    }
}
