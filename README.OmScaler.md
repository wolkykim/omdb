OmScaler
========

OmScaler is a partner program that makes larger OmDB cluster with multiple OmDB backends.

* Extends OmDB to larger scale.
* Sits in front of OmDBs.
* Maps multiple OmDB nodes into linear namespace. (No Consistent-Hash rings)
* Provides data availability by data replication,
* Provides Quorum-based data integrity.

### OmScaler's VS. Cassandra's Consistent-Hash model

The problems of Consistent-Hash based NoSQL like Cassandra.
* It doens't scale well due to the nature of required re-balancing when number of nodes changes. In large scale, server always dies.
* Consistent-hash model is hard to handle hot-spot due to the nature of the algorithm.
* Performance lag grows as cluster grows.
* If serious sever failure happens, it's likely to result global service failure.
* Separate 2nd index is needed to iterate keys.

What about OmScaler's linear partitioning?
* OmScaler uses partition map that maps range of namespace to each set of OmDBs.
* Yeah, there would be some overhead on partition management.
* But much lesser data re-balancing and linear performance increase as cluster grows. So it can scale out.
* Isolated failure. Even if serious of server failure happens, only that range of data gets out of service.
* No index is needed for key iteration.

Note) OmScaler is not publicly available, please contact the author.
