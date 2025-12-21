// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title TowerToken (TOWER) - SECURITY HARDENED
 * @dev Governance/reward token with 0.5% transfer fee
 * @notice Implements: Pausable, ReentrancyGuard, AccessControl, Blacklist, Transfer Fee
 */
contract TowerToken is ERC20, ERC20Burnable, AccessControl, Pausable, ReentrancyGuard {
    // Roles
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");
    bytes32 public constant BLACKLISTER_ROLE = keccak256("BLACKLISTER_ROLE");
    
    // Tokenomics
    uint256 public constant MAX_SUPPLY = 1_000_000 * 10**18; // 1 million TOWER
    
    // Transfer fee (0.5%)
    uint256 public constant TRANSFER_FEE_PERCENT = 50; // 0.5% (50/10000)
    uint256 public constant FEE_DENOMINATOR = 10000;
    address public constant FEE_RECIPIENT = 0xd210925940D236951F0Bfa632A89d76fA1C8883d;
    
    // Platform addresses
    address public platformTreasury;
    address public buybackWallet;
    
    // Per-minter limits
    mapping(address => uint256) public minterDailyLimit;
    mapping(address => mapping(uint256 => uint256)) public minterDailyMinted;
    
    // Blacklist for malicious addresses
    mapping(address => bool) public blacklisted;
    
    // Circuit breaker
    uint256 public dailyVolumeLimit = 10_000_000 * 10**18; // 10M TOWER/day
    mapping(uint256 => uint256) public dailyVolume;
    
    // Events
    event TreasuryUpdated(address indexed oldTreasury, address indexed newTreasury);
    event BuybackWalletUpdated(address indexed oldWallet, address indexed newWallet);
    event MinterAuthorized(address indexed minter, uint256 dailyLimit);
    event MinterRevoked(address indexed minter);
    event Blacklisted(address indexed account, string reason);
    event Unblacklisted(address indexed account);
    event EmergencyPause(address indexed by, string reason);
    // Circuit breaker
    event CircuitBreakerTriggered(uint256 volume, uint256 limit);
    
    // Modifier for blacklist check
    modifier notBlacklisted(address account) {
        require(!blacklisted[account], "Address is blacklisted");
        _;
    }
    
    constructor(
        address _treasury,
        address _buybackWallet,
        address _icoWallet,
        address _devWallet,
        address _marketingWallet,
        address _liquidityWallet
    ) ERC20("Tower Token", "TOWER") {
        require(_treasury != address(0), "Invalid treasury");
        require(_buybackWallet != address(0), "Invalid buyback wallet");
        
        platformTreasury = _treasury;
        buybackWallet = _buybackWallet;
        
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(PAUSER_ROLE, msg.sender);
        _grantRole(BLACKLISTER_ROLE, msg.sender);
        
        // Initial distribution (1M total)
        _mint(_icoWallet, 200_000 * 10**18);         // 20% ICO (200K)
        _mint(_devWallet, 100_000 * 10**18);         // 10% Dev Team (100K)
        _mint(_marketingWallet, 50_000 * 10**18);    // 5% Marketing (50K)
        _mint(_liquidityWallet, 650_000 * 10**18);   // 65% Liquidity (650K)
    }
    
    /**
     * @dev Mint tokens with security checks
     */
    function mint(address to, uint256 amount) 
        external 
        onlyRole(MINTER_ROLE)
        nonReentrant
        whenNotPaused
        notBlacklisted(to)
    {
        require(to != address(0), "Invalid recipient");
        require(amount > 0, "Amount must be > 0");
        require(totalSupply() + amount <= MAX_SUPPLY, "Exceeds max supply");
        
        uint256 today = block.timestamp / 1 days;
        
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
        
        dailyVolume[today] += amount;
        _mint(to, amount);
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
     * @dev Revoke minter
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
     * @dev Add address to blacklist
     */
    function addToBlacklist(address account, string calldata reason) 
        external 
        onlyRole(BLACKLISTER_ROLE) 
    {
        require(account != address(0), "Invalid address");
        require(!hasRole(DEFAULT_ADMIN_ROLE, account), "Cannot blacklist admin");
        blacklisted[account] = true;
        emit Blacklisted(account, reason);
    }
    
    /**
     * @dev Remove address from blacklist
     */
    function removeFromBlacklist(address account) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        blacklisted[account] = false;
        emit Unblacklisted(account);
    }
    
    /**
     * @dev Emergency pause
     */
    function pause(string calldata reason) 
        external 
        onlyRole(PAUSER_ROLE) 
    {
        _pause();
        emit EmergencyPause(msg.sender, reason);
    }
    
    /**
     * @dev Unpause
     */
    function unpause() 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        _unpause();
    }
    
    /**
     * @dev Update treasury
     */
    function updateTreasury(address newTreasury) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        require(newTreasury != address(0), "Invalid address");
        address oldTreasury = platformTreasury;
        platformTreasury = newTreasury;
        emit TreasuryUpdated(oldTreasury, newTreasury);
    }
    
    /**
     * @dev Update buyback wallet
     */
    function updateBuybackWallet(address newWallet) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        require(newWallet != address(0), "Invalid address");
        address oldWallet = buybackWallet;
        buybackWallet = newWallet;
        emit BuybackWalletUpdated(oldWallet, newWallet);
    }
    
    /**
     * @dev Override transfer to add security checks and 0.5% fee
     */
    function _update(address from, address to, uint256 value)
        internal
        override
        whenNotPaused
        notBlacklisted(from)
        notBlacklisted(to)
    {
        // Skip fee for minting/burning
        if (from == address(0) || to == address(0)) {
            super._update(from, to, value);
            return;
        }
        
        // Skip fee for fee recipient
        if (from == FEE_RECIPIENT || to == FEE_RECIPIENT) {
            super._update(from, to, value);
            return;
        }
        
        // Calculate 0.5% fee
        uint256 fee = (value * TRANSFER_FEE_PERCENT) / FEE_DENOMINATOR;
        uint256 amountAfterFee = value - fee;
        
        // Transfer fee to recipient
        if (fee > 0) {
            super._update(from, FEE_RECIPIENT, fee);
        }
        
        // Transfer remaining to recipient
        super._update(from, to, amountAfterFee);
        
        emit TransferWithFee(from, to, value, fee);
    }
    
    event TransferWithFee(
        address indexed from,
        address indexed to,
        uint256 amount,
        uint256 fee
    );
    
    /**
     * @dev Emergency token recovery (only non-TOWER tokens)
     */
    function recoverTokens(address tokenAddress, uint256 amount) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        require(tokenAddress != address(this), "Cannot recover TOWER");
        IERC20(tokenAddress).transfer(msg.sender, amount);
    }
}
