# GoBlockShare
[![GoDoc](https://godoc.org/github.com/denverquane/GoBlockShare?status.png)](https://godoc.org/github.com/denverquane/GoBlockShare)
[![Build Status](https://travis-ci.org/denverquane/GoBlockShare.svg?branch=master)](https://travis-ci.org/denverquane/GoBlockShare)

This app seeks to demonstrate using blockchain tech. for a chat application like Slack or Discord.

By using the blockchain, no user can delete, modify, or falsify chat records and exchanges without other users
being aware of the change. This not ensures integrity of the chat, but allows for interesting functionality 
regarding "rewinding" or rollback of chat engagements.

## Goals:
- [ ] Proof of Work for posting messages/blocks (prevent spam/abuse)
  - [X] Basic difficulty/cryptographic proof validation
  - [ ] Scaling difficulty of blocks as the ~~chain~~ userbase grows
  - [ ] Rewards for propagating the chain, even if not posting
    - Reputation as a global "currency"
    - Individual "tokens" for posting on certain channels -> Have the token = able to post, using REP as the "Gas"
        - Blockchain recognizes when an address is issued a token for a channel, and when it is exchanged/traded away
    - See Bitcoin whitepaper for inspiration: [Bitcoin.org](https://bitcoin.org/bitcoin.pdf)
- [ ] Node discovery
  - [ ] Save chain/channel to disk -> especially for channels with low usercounts w/ possibility of perm. loss
  - [ ] Ability to run app as a node registry/lookup
  - [ ] Active central registry for node lookup
  - [ ] Automatically propagate chain changes to other nodes
  - [ ] Consensus algo. for chain conflicts (explore merging non-conflicting message types/re-generation of transaction)
  - [X] Dockerize app for simplified multi-node testing on a single physical machine
- [ ] Author/poster validation (login validation)
  - [X] Basic authentication
  - [X] ECDSA Public/Private key generation
  - [ ] Tracking addresses across addresses (look how BTC handles "change") -> "Wallet?"
- [ ] Varied communication methods/formats
  - [ ] Private messaging (PGP?) -> Other nodes can carry messages, still
  - [ ] Channel permissions/usergroups? (Only admins can add to private channels...)
    - Private channels can only be propagated by users/admins of that channel...
    - Ensure no channels exist of the same name, across ALL nodes (expensive...)
  - [ ] Different message types (files?)
- [X] Basic JS Frontend for viewing the blockchain in realtime
- [ ] JS Frontend for posting, deleting, editing, etc. messages and transactions
- [ ] Electron integration for running the web app as a native desktop app

## Project Structure:
This app is comprised of two distinct parts as of 4/23/18.
These parts serve as the frontend and backend services for the overarching application.

- #### GO Backend
In the aptly-named "Go" directory are all the [GoLang](https://golang.org/) source files, which run the backend
application that handles blockchain operations, including peer discovery/chain propagation, author and user validation,
and rudimentary ["Proof of Work"](https://en.wikipedia.org/wiki/Proof-of-work_system) calculations to ensure users don't
overwhelm blockchain peers with transactions.

- #### Ts/React Frontend
See [ReactBlockShare](https://github.com/denverquane/ReactBlockShare)