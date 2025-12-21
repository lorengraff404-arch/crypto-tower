// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/token/ERC1155/ERC1155.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC1155/extensions/ERC1155Supply.sol";

/**
 * @title ItemNFT
 * @dev ERC-1155 semi-fungible tokens for items
 * Supports: weapons, armor, accessories, consumables, materials
 */
contract ItemNFT is ERC1155, ERC1155Supply, Ownable {
    // Item metadata
    struct ItemData {
        string itemType;     // WEAPON, ARMOR, ACCESSORY, CONSUMABLE, MATERIAL
        string name;
        string rarity;       // SSS, SS, S, A, B, C
        bool isConsumable;
        uint256 maxSupply;   // 0 = unlimited
        uint256 createdAt;
    }
    
    mapping(uint256 => ItemData) public items;
    mapping(uint256 => string) private _tokenURIs;
    
    // Authorized minters
    mapping(address => bool) public authorizedMinters;
    
    uint256 private _nextItemId;
    
    // Events
    event ItemCreated(uint256 indexed itemId, string itemType, string name, string rarity);
    event MinterAuthorized(address indexed minter);
    event MinterRevoked(address indexed minter);
    
    constructor(
        string memory uri
    ) ERC1155(uri) Ownable(msg.sender) {}
    
    /**
     * @dev Create new item type
     */
    function createItem(
        string memory itemType,
        string memory name,
        string memory rarity,
        bool isConsumable,
        uint256 maxSupply,
        string memory tokenURI
    ) external onlyOwner returns (uint256) {
        uint256 itemId = _nextItemId++;
        
        items[itemId] = ItemData({
            itemType: itemType,
            name: name,
            rarity: rarity,
            isConsumable: isConsumable,
            maxSupply: maxSupply,
            createdAt: block.timestamp
        });
        
        _tokenURIs[itemId] = tokenURI;
        
        emit ItemCreated(itemId, itemType, name, rarity);
        return itemId;
    }
    
    /**
     * @dev Mint items (authorized backend only)
     */
    function mint(
        address to,
        uint256 id,
        uint256 amount
    ) external {
        require(authorizedMinters[msg.sender], "Not authorized");
        require(items[id].createdAt > 0, "Item doesn't exist");
        
        if (items[id].maxSupply > 0) {
            require(
                totalSupply(id) + amount <= items[id].maxSupply,
                "Exceeds max supply"
            );
        }
        
        _mint(to, id, amount, "");
    }
    
    /**
     * @dev Batch mint multiple items
     */
    function mintBatch(
        address to,
        uint256[] memory ids,
        uint256[] memory amounts
    ) external {
        require(authorizedMinters[msg.sender], "Not authorized");
        
        for (uint256 i = 0; i < ids.length; i++) {
            require(items[ids[i]].createdAt > 0, "Item doesn't exist");
            
            if (items[ids[i]].maxSupply > 0) {
                require(
                    totalSupply(ids[i]) + amounts[i] <= items[ids[i]].maxSupply,
                    "Exceeds max supply"
                );
            }
        }
        
        _mintBatch(to, ids, amounts, "");
    }
    
    /**
     * @dev Burn items (for crafting, etc)
     */
    function burn(
        address from,
        uint256 id,
        uint256 amount
    ) external {
        require(
            from == msg.sender || isApprovedForAll(from, msg.sender),
            "Not authorized"
        );
        _burn(from, id, amount);
    }
    
    /**
     * @dev Authorize minter
     */
    function authorizeMinter(address minter) external onlyOwner {
        authorizedMinters[minter] = true;
        emit MinterAuthorized(minter);
    }
    
    /**
     * @dev Revoke minter
     */
    function revokeMinter(address minter) external onlyOwner {
        authorizedMinters[minter] = false;
        emit MinterRevoked(minter);
    }
    
    /**
     * @dev Get item data
     */
    function getItemData(uint256 id) external view returns (ItemData memory) {
        require(items[id].createdAt > 0, "Item doesn't exist");
        return items[id];
    }
    
    /**
     * @dev Get token URI
     */
    function uri(uint256 tokenId) public view override returns (string memory) {
        return _tokenURIs[tokenId];
    }
    
    /**
     * @dev Update token URI
     */
    function setURI(uint256 tokenId, string memory tokenURI) external onlyOwner {
        _tokenURIs[tokenId] = tokenURI;
    }
    
    // Required override
    function _update(address from, address to, uint256[] memory ids, uint256[] memory values)
        internal
        override(ERC1155, ERC1155Supply)
    {
        super._update(from, to, ids, values);
    }
}
