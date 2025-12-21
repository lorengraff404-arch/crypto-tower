// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/token/common/ERC2981.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title CharacterNFT - SECURITY HARDENED
 * @dev ERC-721 with royalties, pausable, and reentrancy protection
 * @notice Implements: ERC2981, Pausable, ReentrancyGuard, AccessControl
 */
contract CharacterNFT is 
    ERC721, 
    ERC721Enumerable, 
    ERC721URIStorage, 
    ERC2981,
    AccessControl, 
    Pausable, 
    ReentrancyGuard 
{
    // Roles
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");
    
    uint256 private _nextTokenId;
    uint256 public constant MAX_SUPPLY = 1_000_000; // 1M max NFTs
    
    // Character metadata
    struct CharacterData {
        uint256 gameCharacterId;
        string characterType;
        string element;
        string rarity;
        uint256 level;
        uint256 mintedAt;
    }
    
    mapping(uint256 => CharacterData) public characters;
    
    // Per-minter daily limits
    mapping(address => uint256) public minterDailyLimit;
    mapping(address => mapping(uint256 => uint256)) public minterDailyMinted;
    
    // Base URI
    string private _baseTokenURI;
    
    // Events
    event CharacterMinted(
        uint256 indexed tokenId,
        address indexed owner,
        uint256 gameCharacterId,
        string characterType,
        string rarity
    );
    event MinterAuthorized(address indexed minter, uint256 dailyLimit);
    event MinterRevoked(address indexed minter);
    event BaseURIUpdated(string newBaseURI);
    event RoyaltyUpdated(address indexed receiver, uint96 feeNumerator);
    event EmergencyPause(address indexed by, string reason);
    
    constructor(
        string memory baseURI
    ) ERC721("Crypto TD Character", "CTDC") {
        _baseTokenURI = baseURI;
        
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(PAUSER_ROLE, msg.sender);
        
        // Set default royalty: 2.5% to contract deployer
        _setDefaultRoyalty(msg.sender, 250); // 250 = 2.5%
    }
    
    /**
     * @dev Mint character NFT with security checks
     */
    function mintCharacter(
        address to,
        uint256 gameCharacterId,
        string memory characterType,
        string memory element,
        string memory rarity,
        uint256 level,
        string memory tokenURI
    ) external onlyRole(MINTER_ROLE) nonReentrant whenNotPaused returns (uint256) {
        require(to != address(0), "Invalid address");
        require(_nextTokenId < MAX_SUPPLY, "Max supply reached");
        require(bytes(characterType).length > 0, "Invalid character type");
        require(bytes(rarity).length > 0, "Invalid rarity");
        
        uint256 today = block.timestamp / 1 days;
        
        // Check per-minter daily limit
        uint256 minterLimit = minterDailyLimit[msg.sender];
        if (minterLimit > 0) {
            require(
                minterDailyMinted[msg.sender][today] < minterLimit,
                "Minter daily limit exceeded"
            );
            minterDailyMinted[msg.sender][today]++;
        }
        
        uint256 tokenId = _nextTokenId++;
        _safeMint(to, tokenId);
        _setTokenURI(tokenId, tokenURI);
        
        characters[tokenId] = CharacterData({
            gameCharacterId: gameCharacterId,
            characterType: characterType,
            element: element,
            rarity: rarity,
            level: level,
            mintedAt: block.timestamp
        });
        
        emit CharacterMinted(tokenId, to, gameCharacterId, characterType, rarity);
        return tokenId;
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
     * @dev Update base URI
     */
    function setBaseURI(string memory baseURI) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        _baseTokenURI = baseURI;
        emit BaseURIUpdated(baseURI);
    }
    
    /**
     * @dev Update royalty info
     */
    function setDefaultRoyalty(address receiver, uint96 feeNumerator) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        _setDefaultRoyalty(receiver, feeNumerator);
        emit RoyaltyUpdated(receiver, feeNumerator);
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
     * @dev Get character data
     */
    function getCharacterData(uint256 tokenId) 
        external 
        view 
        returns (CharacterData memory) 
    {
        require(_ownerOf(tokenId) != address(0), "Token doesn't exist");
        return characters[tokenId];
    }
    
    /**
     * @dev Get all tokens owned by address
     */
    function tokensOfOwner(address owner) 
        external 
        view 
        returns (uint256[] memory) 
    {
        uint256 balance = balanceOf(owner);
        uint256[] memory tokens = new uint256[](balance);
        
        for (uint256 i = 0; i < balance; i++) {
            tokens[i] = tokenOfOwnerByIndex(owner, i);
        }
        
        return tokens;
    }
    
    // Required overrides
    function _baseURI() internal view override returns (string memory) {
        return _baseTokenURI;
    }
    
    function _update(address to, uint256 tokenId, address auth)
        internal
        override(ERC721, ERC721Enumerable)
        whenNotPaused
        returns (address)
    {
        return super._update(to, tokenId, auth);
    }
    
    function _increaseBalance(address account, uint128 value)
        internal
        override(ERC721, ERC721Enumerable)
    {
        super._increaseBalance(account, value);
    }
    
    function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }
    
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, ERC721Enumerable, ERC721URIStorage, ERC2981, AccessControl)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}
