# Commitment Provider - Swap Infrastructure for Monero Trading

The commitment provider aims to tackle a problem with Monero Atomic Swaps, they are difficult to develop.

The commitment provider is a secret provider infrastructure that provides secrets for Hash Time Lock Contracts, to use a simple commitment-reveal scheme while trading monero.

The HTLC contracts can be deployed on Bitcoin, Litecoin, Ethereum, Solana or any chain that can support it.

Monero doesn't have smart contracts and developing multisig is complex and requires the users to download a specialized wallet and complex messaging.

Instead of that, the Commitment Provider is a (decentralizable) third party that verifies Monero transactions and reveals secrets in-exchange.

The commitment provider:

1. Never Has access to any funds
2. Never knows where the commitments and secrets are actually used
3. Can't be used to steal funds by third party, if the database is hacked.
4. It's not an escrow and doesn't provide dispute resolution
5. It's more similar to an Oracle that observes Monero payments and reveals secrets to specific users.
6. Discards the monero payment proofs after the secret is revealed.

The commitment provider reveals the pre-image of a commitment to a user that can provide a valid proof of monero transaction. Payment proofs are checked using `check_tx_key TXID TXKEY ADDRESS`

### Flow

* Alice wants to exchange her XMR to Solana with Bob.
* Alice sends Bob her Sol address and they agree on the exchange rate
* Bob contacts the Commitment Provider and requests a Commitment, When requesting the commitment Bob configures the commitment provider to reveal the secret of the commitment for a payment proof of X amount made to his address
* Bob interacts with a smart contract on Solana and deposits the Sol to trade
* The HTLC on Solana unlocks the Sol if the secret is provided by Alice or it refunds it in 1 day back to Bob
* Alice makes the XMR deposit to Bobs address
* Alice then creates a payment proof and uses it to get the secret from the commitment provider 
* Alice then pulls payments from the Smart contract using the secret

The explained flow works on every smart contract chain and Solana was just an example.

### Benefits
1. Using the Commitment Provider requires no Javascript
2. The commitment provider offers a simple flow for trading
3. It's more trustless than an Escrow
4. Trading is always P2P, Alice transfers directly to Bob
5. No local app to run for users. No asb needed to swap

### Considerations
* The users need to have 2 browser windows open. One to use the commitment provider's interface and one to do a swap with a browser wallet on chain
* It can however be used as a JSON API and embedded into applications/websites or served as an IFRAME
* It must stay online, else there can be loss of funds due to the HTLC timelock expiration, but secret recovery mechanisms can be implemented

### Monetizing
The commitment provider is a paid service that works with a `pay per commitment` model

Example: 
* Bob deposits XMR to the commitment provider's address
* Bob uses the Proof of Deposit to get an API key from the Commitment Provider
* API key has a quota
* Bob uses the API Key to create commitments at the provider and when he uses his quota he needs to get a new key
* Example fees could be: 1000 commitments per $0.1 worth of XMR 


### Decentralization
The commitment provider infrastructure can be decentralized, however instead of building a complex network to create it, a simpler approach should be used, it should work as a DAO.

* N amount of people can become commitment providers
* They create an N of M multisig account via Monero CLI and communicate via chat
* The commitment providers configure their instances to pay into the Multisig account
* At the end of the month, commitment providers withdraw their profits from the multisig in consensus
* If a commitment service goes down and this results in loss of funds, reimbursement could be requested via chat
* A provider that goes down too often could be kicked out if it results in loss of funds

This would allow for a simple decentralization where humans are representing each instance and large networks of nodes don't exist. Having 4-5 providers are more than enough.

It would be also possible to implement decentralized storage for secrets via threshold homomorphic encryption which is up for research (golang latigo library can work for this, but I don't want to over complicate things)

## API

The API serves HTML without Javascript and there is a JSON API

//TODO: This is under development


## How commitments are computed

The commitment is a sha256 hash of a 32 byte secret buffer.
The commitment and the secret can be converted string from to a buffer using using Javascript with this:

```
function decodeString(hexString) {
  return Buffer.from(hexString, 'hex');
}

// Encode a byte array to a hex string
function encodeToString(byteArray) {
  return Buffer.from(byteArray).toString('hex');
}

```

Following this will make sure you stay consistent with the Go implementation which uses `hex.EncodeToString` and `hex.DecodeString` and `sha256.Sum256`

The corresponding hash in javascript is

`const hashBuffer = await crypto.subtle.digest("SHA-256", buff);`

So you can recreate the hash provided by the commitment provider in Javascript easily. This should be compatible with all sha256 hashes provided by on-chain contracts too!