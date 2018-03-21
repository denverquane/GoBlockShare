# GoBlockChat
[![GoDoc](https://godoc.org/github.com/denverquane/GoBlockChat?status.png)](https://godoc.org/github.com/denverquane/GoBlockChat)

This app seeks to demonstrate using blockchain tech. for a chat application like Slack or Discord.

By using the blockchain, no user can delete, modify, or falsify chat records and exchanges without other users
being aware of the change. This not ensures integrity of the chat, but allows for interesting functionality 
regarding "rewinding" or rollback of chat engagements.

## Goals:
- [ ] Proof of Work for posting messages/blocks (prevent spam/abuse)
  - [X] Basic difficulty/cryptographic proof validation
  - [ ] Scaling difficulty of blocks as the chain grows
  - [ ] Rewards for propagating the chain (?)
- [ ] Automatic discovery of other Nodes, and auto propagation of the blockchain itself
- [X] Author/poster validation (login validation)
- [X] Basic JS Frontend for viewing the blockchain in realtime
- [ ] JS Frontend for posting, deleting, editing, etc. messages and transactions
- [ ] GO app to interact with the chain, without the Webapp (?)
