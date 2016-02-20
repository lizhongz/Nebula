/*
Package gossip implements a pull-based protocol for maitaining a membership list
for each Nebula node. Threrefore, nodes know the existence of the others.

When a node join or leave a Nebula cluster, the information is populated through
Gossip protocol. A node also gossips its heartbeat periodly, so that the others
know it is alive.
*/
package gossip
