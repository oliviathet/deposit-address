# Demo

This demo shows how to generate and verify proof of ownership for both Lombard hot wallet addresses and deposit addresses. The hot wallet addresses are Taproot addresses (tb1p..., bc1p...) and the deposit addresses are P2WPKH Segwit addresses(tb1q..., bc1q...). For each address, there is 2 steps:

1. Generate the proof. The proof is a transaction signed by your key with an arbitrarily chosen "challenge" string in the `OP_RETURN` script. The transaction transfers zero value from a set of invalid UTXOs and pays a zero fee, so it cannot be mined.
2. Verify the proof. Since you have the key's address, the proof (the signature), and the challenge string, you can pass these parameters to verify that the signature was made by that key and contains that challenge in the `OP_RETURN` script.

## Generate the proof

If you do not have keys or a CubeSigner account, follow the [Keys Setup Guide](./generate-keys-setup.md). Now that you have your keys, you can generate proofs:

### 1. Configure

Set the following values in your `golang/config.yaml`

| Config          | Description                             | Example            |
|-----------------|-----------------------------------------|--------------------|
| demo.address    | Bitcoin address to generate proof for   | `tb1...`, `bc1...` |
| demo.challenge  | Challenge message to sign               | `TEST`             |
| cubist.session  | Base64 encoded CubeSigner session token |                    |
| cubist.key-id   | Identifier for the signing key          | `Key#Btc...`       |
| cubist.role-id  | Identifier for the authorization role   | `Role#...`         |


If `demo.address` is a deposit address, it will need additional fields, which you can get via the `deposit_metadata` for `demo.address` at https://mainnet.prod.lombard.finance/api/v1/address.

| Config                      | Description                   | Example            |
|-----------------------------|-------------------------------------------|--------------------|
| demo.deposit.to-address     | Deposit address' destination EVM address  | `0x6c5e839bde85b6381f41e8e374797457c68e630b` |
| demo.deposit.to-blockchain  | Deposit address' destination chain        | `DESTINATION_BLOCKCHAIN_ETHEREUM`            |
| demo.deposit.referral       | Deposit address' referral                 | `lombard`      |
| demo.deposit.nonce          | Deposit address' nonce                    | `0`     |


### 2. Sign the proof 
```bash
cd golang
go run demo/cmd/prove/main.go
```


This prints the signature of the proof. You will need this signature for proof verification.


## Verify the proof

To verify proof of ownership, you will need the address and signature from the proof generation. You will also need the challenge; if you didn't specifically pass a challenge string, the default challenge string is "LOMBARD PROOF OF OWNERSHIP".


### 1. Configure
Set the following values in `golang/config.yaml`

| Config          | Description                             | Example            |
|-----------------|-----------------------------------------|--------------------|
| demo.address    | Bitcoin address to verify proof for     | `tb1...`, `bc1...` |
| demo.challenge  | Challenge message to verify for         | `TEST`             |
| demo.signature  | Bitcoin signature to verify proof for   |                    |


### 2. Verify the proof signature
```
cd golang
go run demo/cmd/verify/main.go
```

If the proof is successfully verified, the output will look like
```
LOMBARD PROOF OF OWNERSHIP
Successfully verified proof
```

## How verification works

As you can see in the demo.go code, we first convert the `signature` (the raw deserialized signature) string to a `wire.MsgTx` struct. Then we set the `address`, `signature`, and `challenge` in the `ProofData` struct. We call `Verify()` on this struct which does the following:

1. Determine the address type.
2. Parses the signature and public key.
3. Converts the public key to a tweaked Taproot/Segwit address.
4. Checks that address matches the expected address.
5. Verifies the signature.
6. Checks that the challenge in the `OP_RETURN` script matches the expected challenge.
7. Returns true and no error if all these checks pass.

