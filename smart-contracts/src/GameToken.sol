// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title GameToken (GTK) - SECURITY HARDENED
 * @dev In-game utility token for opBNB with comprehensive security
 * @notice Implements: Pausable, ReentrancyGuard, AccessControl, Rate Limiting
 */
contract GameToken is ERC20, ERC20Burnable, AccessControl, Pausable, ReentrancyGuard {
    // Roles
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");
    
    // Tokenomics - CORRECTED (100M GTK : 1M TOWER = 100:1 ratio)
    uint256 public constant MAX_SUPPLY = 100_000_000 * 10**18; // 100 million GTK
    uint256 public dailyEmissionCap = 1_000_000 * 10**18; // 1M GTK per day
    
    // Daily emission tracking
    mapping(uint256 => uint256) public dailyEmissions;
    
    // Per-minter daily limits
    mapping(address => uint256) public minterDailyLimit;
    mapping(address => mapping(uint256 => uint256)) public minterDailyMinted;
    
    // Circuit breaker
    uint256 public dailyVolumeLimit = 5_000_000 * 10**18; // 5M GTK/day (5% of supply)
    mapping(uint256 => uint256) public dailyVolume;
    
    // Events
    event MinterAuthorized(address indexed minter, uint256 dailyLimit);
    event MinterRevoked(address indexed minter);
    event DailyCapUpdated(uint256 newCap);
    event EmergencyPause(address indexed by, string reason);
    event EmergencyUnpause(address indexed by);
    event CircuitBreakerTriggered(uint256 volume, uint256 limit);
    
    constructor() ERC20("Game Token", "GTK") {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(PAUSER_ROLE, msg.sender);
        
        // Initial treasury mint (10% of max supply = 10M GTK)
        _mint(msg.sender, 10_000_000 * 10**18);
    }
    
    /**
     * @dev Mint tokens with comprehensive security checks
     * @param to Recipient address
     * @param amount Amount to mint
     */
    function mint(address to, uint256 amount) 
        external 
        onlyRole(MINTER_ROLE)
        nonReentrant
        whenNotPaused
    {
        require(to != address(0), "Invalid recipient");
        require(amount > 0, "Amount must be > 0");
        require(totalSupply() + amount <= MAX_SUPPLY, "Exceeds max supply");
        
        uint256 today = block.timestamp / 1 days;
        
        // Check global daily emission cap
        require(
            dailyEmissions[today] + amount <= dailyEmissionCap,
            "Daily emission cap exceeded"
        );
        
        // Check per-minter daily limit
        uint256 minterLimit = minterDailyLimit[msg.sender];
        if (minterLimit > 0) {
            require(
                minterDailyMinted[msg.sender][today] + amount <= minterLimit,
                "Minter daily limit exceeded"
            );
            minterDailyMinted[msg.sender][today] += amount;
        }
        
        // Check circuit breaker
        require(
            dailyVolume[today] + amount <= dailyVolumeLimit,
            "Circuit breaker: daily volume limit exceeded"
        );
        
        dailyEmissions[today] += amount;
        dailyVolume[today] += amount;
        
        _mint(to, amount);
        
        emit Transfer(address(0), to, amount);
    }
    
    /**
     * @dev Authorize minter with daily limit
     */
    function authorizeMinter(address minter, uint256 dailyLimit) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        require(minter != address(0), "Invalid minter");
        _grantRole(MINTER_ROLE, minter);
        minterDailyLimit[minter] = dailyLimit;
        emit MinterAuthorized(minter, dailyLimit);
    }
    
    /**
     * @dev Revoke minter authorization
     */
    function revokeMinter(address minter) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        _revokeRole(MINTER_ROLE, minter);
        minterDailyLimit[minter] = 0;
        emit MinterRevoked(minter);
    }
    
    /**
     * @dev Update daily emission cap
     */
    function updateDailyCap(uint256 newCap) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        require(newCap <= MAX_SUPPLY / 10, "Cap too high");
        dailyEmissionCap = newCap;
        emit DailyCapUpdated(newCap);
    }
    
    /**
     * @dev Emergency pause (circuit breaker)
     */
    function pause(string calldata reason) 
        external 
        onlyRole(PAUSER_ROLE) 
    {
        _pause();
        emit EmergencyPause(msg.sender, reason);
    }
    
    /**
     * @dev Unpause after emergency
     */
    function unpause() 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        _unpause();
        emit EmergencyUnpause(msg.sender);
    }
    
    /**
     * @dev Get today's remaining mint capacity
     */
    function getRemainingDailyCapacity() external view returns (uint256) {
        uint256 today = block.timestamp / 1 days;
        uint256 emitted = dailyEmissions[today];
        return emitted >= dailyEmissionCap ? 0 : dailyEmissionCap - emitted;
    }
    
    /**
     * @dev Get minter's remaining daily capacity
     */
    function getMinterRemainingCapacity(address minter) external view returns (uint256) {
        uint256 today = block.timestamp / 1 days;
        uint256 limit = minterDailyLimit[minter];
        if (limit == 0) return 0;
        
        uint256 minted = minterDailyMinted[minter][today];
        return minted >= limit ? 0 : limit - minted;
    }
    
    /**
     * @dev Override transfer to add pause check
     */
    function _update(address from, address to, uint256 value)
        internal
        override
        whenNotPaused
    {
        super._update(from, to, value);
    }
}
